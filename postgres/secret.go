package postgres

import (
	utils "github.com/iboware/postgresql-operator/utils"

	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NewPostgresSecret creates secrets for StatefulSet
func NewPostgresSecret(cr *databasev1alpha1.PostgreSQL) *v1.Secret {
	labels := utils.NewLabels("postgres", cr.ObjectMeta.Name, "postgres")
	encodedpgPassword := cr.Spec.PostgresPassword
	encodedrepMgrPassword := cr.Spec.RepMGRPassword

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name + "-postgresql",
			Labels:    labels,
			Namespace: cr.Spec.Namespace,
		},
		Type: v1.SecretTypeOpaque,
		StringData: map[string]string{
			"postgresql-password": encodedpgPassword,
			"repmgr-password":     encodedrepMgrPassword,
			"admin-password":      encodedpgPassword,
		},
	}
	return secret
}
