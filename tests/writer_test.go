package tests

import(
	"testing"
	"syscall"
	"path/filepath"
	"io"

	"os"
	"log"
	"fmt"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"

 .	"github.com/byteball/go-byteballcore/types"

	"github.com/byteball/go-byteballcore/conf"
	"github.com/byteball/go-byteballcore/db"
	"github.com/byteball/go-byteballcore/storage"
	"github.com/byteball/go-byteballcore/writer"
)

type(
	DBConnT		= db.DBConnT

	ObjJointT	= writer.ObjJointT
	ObjValidationStateT = writer.ObjValidationStateT

	UnitObjectT	= writer.UnitObjectT
)


type(
	oJoVST struct{
		oJ	*ObjJointT
		oVS	*ObjValidationStateT
	}
)

func getReplayFile(rfn string) []oJoVST {
	fmt.Printf("replay file %s\n", rfn)

	rfc, erf := ioutil.ReadFile(rfn)
	if erf != nil {
		log.Panicf("%s: %s", rfn, erf.Error())
	}

	rfs := strings.Split(string(rfc), "\n")
	fmt.Printf("replay lines %d\n", len(rfs))

	oJoVSs := make([]oJoVST, 0, len(rfs)/3)
	for j:=0; j<len(rfs)-1; j+=3 {
		if rfs[j] != "oJ,oVS" {
			log.Panicf("rfs[%d] %s\n", j, rfs[j])
		}

		oJ := ObjJointT{}

		// [tbd] abstract this as func (*T) Make()
		oJ.Unit.Main_chain_index = MCIndexT_Null

		eoj := json.Unmarshal([]byte(rfs[j+1]), &oJ)
		if eoj != nil {
			log.Panicf("rfs[%d] %s\n", j, eoj.Error())
		}

		oVS := ObjValidationStateT{}

		// [tbd] abstract this as func (*T) Make()
		oVS.Max_parent_limci = MCIndexT_Null
		oVS.Last_ball_mci = MCIndexT_Null
		oVS.Max_known_mci = MCIndexT_Null

		eojs := json.Unmarshal([]byte(rfs[j+2]), &oVS)
		if eojs != nil {
			log.Panicf("rfs[%d] %s\n", j, eojs.Error())
		}

		oJoVSs = append(oJoVSs, oJoVST{ oJ: &oJ, oVS: &oVS })
	}

	return oJoVSs
}

func dbStats() {
	conn := db.TakeConnectionFromPool_sync()
	defer conn.Release()
	conn.ShowPrepared()
}

func showStats() {
	dbStats()
}

func redirectOutputToFile(lfn string) {
	fmt.Printf("log file %s\n", lfn)

	// [fyi] save stdout/stderr hanldes
	syscall.Dup2(1, 64+1)
//	syscall.Dup2(2, 64+2)

	lfl, _ := os.OpenFile(lfn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0755)
	syscall.Dup2(int(lfl.Fd()), 1)
//	syscall.Dup2(int(lfl.Fd()), 2)
}

func redirectOutputCancel() {
	// [fyi] restore stdout/stderr hanldes
	syscall.Dup2(64+1, 1)
//	syscall.Dup2(64+2, 2)
}


const(
	dbfnPrefix = "byteball.sqlite"
)

func setupReplay(srcDir string) ([]oJoVST) {
	adn := conf.AppDirName()
	fmt.Printf("AppDirName %s\n", adn)

	if _, err := os.Stat(adn); os.IsNotExist(err) {
		err := os.MkdirAll(adn, os.ModePerm)
		if err != nil {
			//log.Fatalf("error creating %s: %s", adn, err.Error())
			log.Fatalf("%s", err.Error())
		}
	}
	if fif, err := os.Stat(adn); err != nil || ! fif.IsDir() {
		//log.Fatalf("error accessing %s: %s", adn, err.Error())
		log.Fatalf("%s", err.Error())
	}

	cwd, _ := os.Getwd()
	fmt.Printf("cwd %s\n", cwd)

	srd := filepath.Join(cwd, srcDir)

	sfs, err := ioutil.ReadDir(srd)
	if err != nil {
		//log.Fatalf("error accessing %s: %s", srd, err.Error())
		log.Fatalf("%s", err.Error())
	}

	for _, sfi := range sfs {
		dbfn := sfi.Name()
		if strings.HasPrefix(dbfn, dbfnPrefix) {
			// [tbd] don't copy again
			fmt.Printf("copying %s\n", dbfn)
			ifn := filepath.Join(srd, dbfn)
			ofn := filepath.Join(adn, dbfn)
			copyFile(ifn, ofn)
		}
	}

	rfn := filepath.Join(srd, "saveJoint.replay.lst")
	oJoVSs := getReplayFile(rfn)

	return oJoVSs
}


func copyFile(src, dst string) (err error) {
	inf, err := os.Open(src)
	if err != nil { return }
	defer inf.Close()

	ouf, err := os.Create(dst)
	if err != nil { return }
	defer func() {
		erc := ouf.Close()
		if err == nil { err = erc }
	}()

	if _, err = io.Copy(ouf, inf); err != nil { return }

	err = ouf.Sync()
	return
}



func TestWriter_Replay_Tiny(t *testing.T) {
	writerReplay(t, "tiny")
}

func TestWriter_Replay_Small(t *testing.T) {
	writerReplay(t, "small")
}

var(
	writerReplaysInitiated int = 0
)

func writerReplay(t *testing.T, tag string) {

	if 0 < writerReplaysInitiated {
		t.Fatalf("only one writer replay per test run")
	}
	writerReplaysInitiated++

	oJoVSs := setupReplay("testdata/replay/" + tag)

	tag_ := strings.Title(tag)
	lfn := "/tmp/go-byteball_Writer_Replay_" + tag_ + ".log"
	redirectOutputToFile(lfn)
	defer redirectOutputCancel()

	db.Init()
	storage.Init()

	for _, oJoVS := range oJoVSs {

		oJ := oJoVS.oJ
		oVS := oJoVS.oVS

		db.TReset()

	        t0 := time.Now()

		writer.SaveJoint_sync(oJ, oVS, func (conn *DBConnT) {
		})

		dt := time.Now().Sub(t0).Nanoseconds()
		fmt.Printf("\n... saveJoint:  %.3f\n", float64(dt) / 1.0e6)

		{{
		tx := db.TExec
		tq := db.TQuery
		tc := db.TCommit
		tr := db.TRollback

		tdb := tx + tq + tc + tr

		fmt.Printf("... dbtimes:  %.3f  %4.1f%%  %8.3f %8.3f %8.3f %8.3f\n",
				float64(tdb) / 1.0e6,
				(float64(tdb) / float64(dt)) * 100.0,
				float64(tx) / 1.0e6,
				float64(tq) / 1.0e6,
				float64(tc) / 1.0e6,
				float64(tr) / 1.0e6)
		}}

		//break
	}

	showStats()

//	redirectOutputCancel()

} // writerReplay
