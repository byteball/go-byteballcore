package object_hash

import(
	"log"
	"fmt"
	"sort"
	"strings"
	"strconv"

	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"

 _core	"nodejs/core"
	"nodejs/console"

 .	"github.com/byteball/go-byteballcore/types"
	"github.com/byteball/go-byteballcore/chash"
)

type(
	refCHash160T	= *CHash160T
	refHashBase64T	= *HashBase64T
	refSHA1HexT	= *SHA1HexT

	refAddressesT	= *AddressesT
	refBallsT	= *BallsT
)


func GetChash160(obj interface{}) CHash160T {
//	panic("[tbd] GetChash160")
	objSS := GetSourceString(obj)
	console.Log("objSS %#v", objSS)
	return chash.CHash160([]byte(objSS))
}

func GetBase64Hash(obj interface{}) refHashBase64T {
//	panic("[tbd] GetBase64Hash")
	objSS := GetSourceString(obj)
	console.Log("objSS %#v", objSS)
	return HashBase64(objSS)
}

func GetUnitContentHash(unit UnitT) int {
	panic("[tbd] GetUnitContentHash")
}

func GetBallHash(unit UnitT, arrParentBalls BallsT, arrSkiplistBalls BallsT, bNonserial bool) refHashBase64T {
	// {
	//	// [fyi] in alphabetical order:
	//	is_nonserial: true,
	//	parent_balls: arrParentBalls,
	//	skiplist_balls: arrSkiplistBalls,
	//	unit: unit,
	// }
	ss := (func () string {
		oss := []string{}
		if bNonserial {
			oss = append(oss, "is_nonserial")
			oss = append(oss, "true")
		}
		if 0 < len(arrParentBalls) {
			oss = append(oss, "parent_balls")
			oss = append(oss, getSSBalls(&arrParentBalls))
		}
		if 0 < len(arrSkiplistBalls) {
			oss = append(oss, "skiplist_balls")
			oss = append(oss, getSSBalls(&arrSkiplistBalls))
		}
		oss = append(oss, "unit")
		oss = append(oss, "s", string(unit))
		return strings.Join(oss, STRING_JOIN_CHAR)
	})()
	hash := HashBase64(ss)
	//console.Log("GetBallHash %#v %#v %#v %#v %#v", unit, arrParentBalls, arrSkiplistBalls, ss, *hash)
	return hash
}

func HashBase64(ss string) refHashBase64T {
	// [fyi] strings are UTF-8 for []byte conversion by default
	bs := []byte(ss)
	hash := sha256.Sum256(bs)
	digest := HashBase64T(base64.StdEncoding.EncodeToString(hash[:]))
	return &digest
}

func SHA1Hex(ss string) refSHA1HexT {
	// [fyi] strings are UTF-8 for []byte conversion by default
	bs := []byte(ss)
	hash := sha1.Sum(bs)
	digest := SHA1HexT(hex.EncodeToString(hash[:]))
	return &digest
}



/**
function getBase64Hash(obj) {
        return crypto.createHash("sha256").update(getSourceString(obj), "utf8").digest("base64");
}
 **/

/**
function getBallHash(unit, arrParentBalls, arrSkiplistBalls, bNonserial) {
        var objBall = {
                unit: unit
        };
        if (arrParentBalls && arrParentBalls.length > 0)
                objBall.parent_balls = arrParentBalls;
        if (arrSkiplistBalls && arrSkiplistBalls.length > 0)
                objBall.skiplist_balls = arrSkiplistBalls;
        if (bNonserial)
                objBall.is_nonserial = true;
        return getBase64Hash(objBall);
}
 **/



func GetSourceString(obj interface{}) string {
	oss := []string{}
	switch v := obj.(type) {

	case bool:
		boolean := "false"
		if v { boolean = "true" }
		oss = append(oss, "b", boolean)

	case int:
		numeric := strconv.Itoa(v)
		oss = append(oss, "n", numeric)

	case float64:
		numeric := fmt.Sprintf("%v", v)
		if iv := int(v); float64(iv) == v {
			numeric = strconv.Itoa(iv)
		}
		oss = append(oss, "n", numeric)

	case string:
		oss = append(oss, "s", v)

	case []interface{}:
		oss = append(oss, "[")
		for _, value := range v {
			oss = append(oss, GetSourceString(value))
		}
		oss = append(oss, "]")

	case AddressT:
		oss = append(oss, "s", string(v))

	case AddressesT:
		oss = append(oss, "[")
		for _, value := range v {
			oss = append(oss, GetSourceString(value))
		}
		oss = append(oss, "]")

	case map[string] interface{}:
		keys := []string{}
		for key := range v {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		for _, key := range keys {
			oss = append(oss, key, GetSourceString(v[key]))
		}

	default:
		var b bool = obj.(bool)
		b = b
		log.Panicf("[tbd] GetSourceString %#v", obj)

	}
	return strings.Join(oss, STRING_JOIN_CHAR)
}

const(
	STRING_JOIN_CHAR	= "\x00"
)

func getSSAddresses(addresses refAddressesT) string {
	if len(*addresses) == 0 {
		_core.Throw("getSSAddresses: empty array")
	}
	oss := []string{}
	oss = append(oss, "[")
	for _, address := range *addresses {
		oss = append(oss, "s", string(address))
	}
	oss = append(oss, "]")
	return strings.Join(oss, STRING_JOIN_CHAR)
}

func getSSBalls(balls refBallsT) string {
	if len(*balls) == 0 {
		_core.Throw("getSSBalls: empty array")
	}
	oss := []string{}
	oss = append(oss, "[")
	for _, ball := range *balls {
		oss = append(oss, "s", string(ball))
	}
	oss = append(oss, "]")
	return strings.Join(oss, STRING_JOIN_CHAR)
}


/***
/*jslint node: true * /
"use strict";

var STRING_JOIN_CHAR = "\x00";

/**
 * Converts the argument into a string by mapping data types to a prefixed string and concatenating all fields together.
 * @param obj the value to be converted into a string
 * @returns {string} the string version of the value
 * /
function getSourceString(obj) {
    var arrComponents = [];
    function extractComponents(variable){
        if (variable === null)
            throw Error("null value in "+JSON.stringify(obj));
        switch (typeof variable){
            case "string":
                arrComponents.push("s", variable);
                break;
            case "number":
                arrComponents.push("n", variable.toString());
                break;
            case "boolean":
                arrComponents.push("b", variable.toString());
                break;
            case "object":
                if (Array.isArray(variable)){
                    if (variable.length === 0)
                        throw Error("empty array in "+JSON.stringify(obj));
                    arrComponents.push('[');
                    for (var i=0; i<variable.length; i++)
                        extractComponents(variable[i]);
                    arrComponents.push(']');
                }
                else{
                    var keys = Object.keys(variable).sort();
                    if (keys.length === 0)
                        throw Error("empty object in "+JSON.stringify(obj));
                    keys.forEach(function(key){
                        if (typeof variable[key] === "undefined")
                            throw Error("undefined at "+key+" of "+JSON.stringify(obj));
                        arrComponents.push(key);
                        extractComponents(variable[key]);
                    });
                }
                break;
            default:
                throw Error("hash: unknown type="+(typeof variable)+" of "+variable+", object: "+JSON.stringify(obj));
        }
    }

    extractComponents(obj);
    return arrComponents.join(STRING_JOIN_CHAR);
}

exports.STRING_JOIN_CHAR = STRING_JOIN_CHAR; // for tests
exports.getSourceString = getSourceString;


 ***/
