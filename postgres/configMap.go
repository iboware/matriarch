package postgres

import (
	"io/ioutil"
	"path/filepath"

	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	utils "github.com/iboware/postgresql-operator/utils"
	"github.com/prometheus/common/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NewConfigMap is a fucntion which creates the configmap for hooks scripts
func NewConfigMap(cr *databasev1alpha1.PostgreSQL) (*v1.ConfigMap, error) {
	labels := utils.NewLabels("postgres", cr.ObjectMeta.Name, "postgres")
	absPath, _ := filepath.Abs("scripts/pre-stop.sh")
	preStop, err := ioutil.ReadFile(absPath)

	if err != nil {
		log.Error(err, "Unable to read file.")
		return nil, err
	}
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name + "-postgresql-hooks-scripts",
			Labels:    labels,
			Namespace: cr.Spec.Namespace,
		},
		Data: map[string]string{
			"pre-stop.sh": string(preStop),
		},
	}

	return configMap, nil
}
