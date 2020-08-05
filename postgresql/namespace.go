package postgresql

import (
	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NewNamespaceForCR creates a namespace for StatefulSet
func NewNamespaceForCR(cr *databasev1alpha1.PostgreSQL) *v1.Namespace {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Spec.Namespace,
		},
	}
	return ns
}
