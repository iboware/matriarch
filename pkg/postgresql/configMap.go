package postgresql

import (
	"io/ioutil"
	"path/filepath"
	postgresqlv1alpha1 "postgres-operator/pkg/apis/postgresql/v1alpha1"
	"postgres-operator/pkg/utils"

	"github.com/prometheus/common/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NewConfigMapForCR is a fucntion which creates the configmap for hooks scripts
func NewConfigMapForCR(cr *postgresqlv1alpha1.PostgreSQL) (*v1.ConfigMap, error) {
	labels := utils.LabelsForPostgreSQL(cr.ObjectMeta.Name)
	absPath, _ := filepath.Abs("../../scripts/pre-stop.sh")
	preStop, err := ioutil.ReadFile(absPath)

	if err != nil {
		log.Error(err, "Unable to read file.")
		return nil, err
	}
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: cr.ObjectMeta.Name + "-postgresql-hooks-scripts", Labels: labels, Namespace: cr.ObjectMeta.Namespace},
		Data: map[string]string{
			"pre-stop.sh": string(preStop),
		},
	}

	return configMap, nil
}
