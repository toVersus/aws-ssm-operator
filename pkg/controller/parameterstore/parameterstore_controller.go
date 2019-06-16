package parameterstore

import (
	"context"
	"fmt"

	errs "github.com/pkg/errors"
	ssmv1alpha1 "github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("parameterstore-controller")

// Add creates a new ParameterStore Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileParameterStore{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("parameterstore-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ParameterStore
	err = c.Watch(&source.Kind{Type: &ssmv1alpha1.ParameterStore{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ssmv1alpha1.ParameterStore{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileParameterStore implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileParameterStore{}

// ReconcileParameterStore reconciles a ParameterStore object
type ReconcileParameterStore struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme

	ssmc *SSMClient
}

// Reconcile reads that state of the cluster for a ParameterStore object and makes changes based on the state read
// and what is in the ParameterStore.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileParameterStore) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ParameterStore")

	// Fetch the ParameterStore instance
	instance := &ssmv1alpha1.ParameterStore{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Secret object
	desired, err := r.newSecretForCR(instance)
	if err != nil {
		return reconcile.Result{}, errs.Wrap(err, "failed to compute secret for cr")
	}

	// Set ParameterStore instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	current := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, current)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Secret", "desired.Namespace", desired.Namespace, "desired.Name", desired.Name)
		err = r.client.Create(context.TODO(), desired)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Secret already exists", "current.Namespace", current.Namespace, "current.Name", current.Name)
	return reconcile.Result{}, nil
}

// newSecretForCR returns a Secret with the same name/namespace as the cr
func (r *ReconcileParameterStore) newSecretForCR(cr *ssmv1alpha1.ParameterStore) (*corev1.Secret, error) {
	labels := map[string]string{
		"app": cr.Name,
	}
	if r.ssmc == nil {
		r.ssmc = newSSMClient(nil)
	}
	ref := cr.Spec.ValueFrom.ParameterStoreRef
	log.Info(fmt.Sprintf("parameterstore name: %s", ref.Name))
	data, err := r.ssmc.SSMParameterValueToSecret(ref)
	if err != nil {
		return nil, errs.Wrap(err, "failed to get json secret as map")
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		StringData: data,
	}, nil
}
