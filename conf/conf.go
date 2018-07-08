package conf

import(
	"os"
	"path"

//	"log"
//	"io/ioutil"
//	"encoding/json"

	db "github.com/byteball/go-byteballcore/db/conf"
)

type(
)

var(
	IsLight 	bool
	Storage		string
)

type Conf struct{
	DB		db.Conf
/**
	Server		server.Conf
	NodeSet		nodeSet.Conf
	//Node		node.Conf
	Accounts	accounts.Conf
	Transactions	transactions.Conf
 **/

}


var(
	confInstance *Conf = &Conf{}
)

func Instance() *Conf {
	return confInstance
}

func AppDirName() string {
//	appName := "byteball-headless"
	appName := "headless-byteball"
	adn := path.Join(os.Getenv("HOME"), ".config", appName)
	return adn
}

func confFileName() string {
//	cfndir := os.Getenv("HOME") + "/" + ".config" + "/" + "byteball-headless"
//	cfn := cfndir + "/" + "conf.json"
	cfn := path.Join(AppDirName(), "conf.json")

	// [tbd] --conf <file>
	// [tbd] env

	return cfn
}

func Init() {
/**
	cfn := confFileName()

	jsonb, errf := ioutil.ReadFile(cfn)
	if errf != nil {
		log.Fatalf("conf: %s", errf.Error())
	}

	errum := json.Unmarshal(jsonb, &confInstance)
	if errum != nil {
		log.Fatalf("conf: %s", errum.Error())
	}
 **/
}
