package archiving

import(
 .	"github.com/byteball/go-byteballcore/types"

	"github.com/byteball/go-byteballcore/db"
)

type(
	DBConnT		= db.DBConnT
	refDBConnT	= *DBConnT

	refJointT	= *JointT

	refAsyncFunctorsT = *AsyncFunctorsT
)

//  (*db.DBConnT, *types.JointT, string, *[]func() error)

func GenerateQueriesToArchiveJoint_sync(conn refDBConnT, objJoint refJointT, s string, arrQueries refAsyncFunctorsT) {
	panic("[tbd] archiving.GenerateQueriesToArchiveJoint")
}
