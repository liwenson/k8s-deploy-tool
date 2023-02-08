package logic

import (
	"NBServer/tools/deploymentgen/internal/jsonToyaml/config"
	"fmt"
	"github.com/ghodss/yaml"
	"log"
	"os"
)

func JsonToYaml(c config.Content) {
	//j := []byte(`{"default":{"namespace":"dev","replicas":1,"resources":{"requests":{"cpu":"20m","memory":"20Mi"},"limits":{"cpu":"100m","memory":"100Mi"}},"volumes":[{"name":"test","hostPath":{"path":"/data"}},{"name":"site-data","persistentVolumeClaim":{"claimName":"my-lamp-site-data"}}]},"services":[{"name":"be-bbs-test-logdemo","image":"reg.test.com/be-ssb-test-logdemo:20220421103226_101","host":"www.test.com","ports":[{"port":8080,"type":"application"}]}]}`)

	jStr := []byte(c.Content)

	yStr, err := yaml.JSONToYAML(jStr)
	if err != nil {
		panic(err)
	}

	if c.Write != "" {
		err = os.WriteFile(c.Write, yStr, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(string(yStr))
	}
}
