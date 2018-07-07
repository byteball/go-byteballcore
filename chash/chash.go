package chash

import(
//	"fmt"
	"strings"
	"strconv"
//	"crypto/sha1"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"encoding/base32"
	"encoding/base64"

 _core	"nodejs/core"
	"nodejs/console"

 .	"github.com/byteball/go-byteballcore/types"

)

const(
	CHASH160_BITS		= 32*5
	CHASH288_BITS		= 48*6

	CHECKSUM_BITS		= 32
)

type(
	OffsetsT	[]int
//	refOffsetsT	= *OffsetsT
)


const PI string = "14159265358979323846264338327950288419716939937510"

var arrRelativeOffsets []string = strings.Split(PI, "")

func checkLength(chash_length int) {
	if chash_length != CHASH160_BITS && chash_length != CHASH288_BITS {
//		throw Error("unsupported c-hash length: "+chash_length);
		_core.Throw("unsupported c-hash length: %d", chash_length);
	}
}

func calcOffsets(chash_length int) OffsetsT {
	checkLength(chash_length)

	arrOffsets := make(OffsetsT, 0, chash_length)
	offset := 0
	index := 0

//	for (var i=0; offset<chash_length; i++){
	for _, strRelativeOffset := range arrRelativeOffsets {
//		var relative_offset = parseInt(arrRelativeOffsets[i]);
		relative_offset, _ := strconv.Atoi(strRelativeOffset)
		if relative_offset == 0 {
			continue
		}
		offset += relative_offset
		if chash_length == CHASH288_BITS {
			offset += 4
		}
		if offset >= chash_length {
			break
		}
		arrOffsets = append(arrOffsets, offset)
		//console.Log("index=%d, offset=%d", index, offset);
		index++;
	}

	if index != CHECKSUM_BITS {
//		throw Error("wrong number of checksum bits");
		_core.Throw("wrong number of checksum bits: %d", index);
	}

	return arrOffsets
}

//var arrOffsets160 = calcOffsets(160);
//var arrOffsets288 = calcOffsets(288);
var arrOffsets160 OffsetsT
var arrOffsets288 OffsetsT

func init() {
	arrOffsets160 = calcOffsets(CHASH160_BITS)
	arrOffsets288 = calcOffsets(CHASH288_BITS)
}

type(
	BitsT		[]int
	refBitsT	*BitsT
)

func separateIntoCleanDataAndChecksum(bin BitsT) (BitsT, BitsT) {
//	var len = bin.length;
	binlen := len(bin)
//	var arrOffsets;
	var arrOffsets OffsetsT
	switch binlen {
//	if (len === 160)
	case CHASH160_BITS:
		arrOffsets = arrOffsets160
//	else if (len === 288)
	case CHASH288_BITS:
		arrOffsets = arrOffsets288
//	else
	default:
//		throw Error("bad length="+len+", bin = "+bin);
		_core.Throw("bad length=%d, bin = %#v", binlen, bin);
	}
	arrFrags := []BitsT{}
	arrChecksumBits := []BitsT{}
	start := 0
//	for (var i=0; i<arrOffsets.length; i++){
	for _, offset := range arrOffsets {
//		arrFrags.push(bin.substring(start, arrOffsets[i]));
		arrFrags = append(arrFrags, bin[start:offset])
//		arrChecksumBits.push(bin.substr(arrOffsets[i], 1));
		arrChecksumBits = append(arrChecksumBits, bin[offset:offset+1])
		start = offset + 1
	}
	// add last frag
	if start < binlen {
//		arrFrags.push(bin.substring(start));
		arrFrags = append(arrFrags, bin[start:])
	}
//	var binCleanData = arrFrags.join("");
	binCleanData := make(BitsT, 0, binlen)
	for _, frag := range arrFrags {
		binCleanData = append(binCleanData, frag...)
	}
//	var binChecksum = arrChecksumBits.join("");
	binChecksum := make(BitsT, 0, CHECKSUM_BITS)
	for _, checksumBit := range arrChecksumBits {
		binChecksum = append(binChecksum, checksumBit...)
	}
//	return {clean_data: binCleanData, checksum: binChecksum};
	return binCleanData, binChecksum
}

func mixChecksumIntoCleanData(binCleanData BitsT, binChecksum BitsT) BitsT {
	if len(binChecksum) != CHECKSUM_BITS {
//		throw Error("bad checksum length");
		_core.Throw("bad checksum length %d", len(binChecksum));
	}
	mixlen := len(binCleanData) + len(binChecksum);
	var arrOffsets OffsetsT
	switch mixlen {
//	if (len === 160)
	case CHASH160_BITS:
		arrOffsets = arrOffsets160
//	else if (len === 288)
	case CHASH288_BITS:
		arrOffsets = arrOffsets288
//	else
	default:
//		throw Error("bad length="+len+", clean data = "+binCleanData+", checksum = "+binChecksum);
		_core.Throw("bad length=%d, clean data = %#v, checksum = %#v", mixlen, binCleanData, binChecksum)
	}
	arrFrags := []BitsT{}
//	var arrChecksumBits = binChecksum.split("");
	start := 0
//	for (var i=0; i<arrOffsets.length; i++){
	for i, offset := range arrOffsets {
		end := offset - i
//		arrFrags.push(binCleanData.substring(start, end));
		arrFrags = append(arrFrags, binCleanData[start:end])
//		arrFrags.push(arrChecksumBits[i]);
		arrFrags = append(arrFrags, binChecksum[i:i+1])
		start = end
	}
	// add last frag
	if start < len(binCleanData) {
//		arrFrags.push(binCleanData.substring(start));
		arrFrags = append(arrFrags, binCleanData[start:])
	}
	//console.Log("arrFrags %#v", arrFrags)
//	return arrFrags.join("");
	binMix := make(BitsT, 0, mixlen)
	for _, frag := range arrFrags {
		binMix = append(binMix, frag...)
	}
	return binMix
}

func buffer2bin(buf []byte) BitsT {
//	var bytes = [];
	bits := make(BitsT, 0, len(buf)*8)
//	for (var i=0; i<buf.length; i++){
	for _, buf_i := range buf {
//		var bin = buf[i].toString(2);
//		if (bin.length < 8) // pad with zeros
//			bin = zeroString.substring(bin.length, 8) + bin;
//		bytes.push(bin);
		for j:=0; j<8; j++ {
			bit := int((buf_i >> uint(7-j)) & 1)
			bits = append(bits, bit)
		}
	}
//	return bytes.join("");
	return bits
}

func bin2buffer(bin BitsT) []byte {
	buflen := len(bin)/8
	buf := make([]byte, 0, buflen)
	for i:=0; i<buflen; i++ {
//		buf[i] = parseInt(bin.substr(i*8, 8), 2);
		byte_i := 0
		for j:=0; j<8; j++ {
			byte_i <<= 1
			byte_i |= (bin[i*8+j] & 1)
		}
		buf = append(buf, byte(byte_i))
	}
	return buf
}

func checksum(clean_data []byte) []byte {
//	var full_checksum = crypto.createHash("sha256").update(clean_data).digest();
	fc := sha256.Sum256(clean_data)
	//console.log(full_checksum);
//	var checksum = new Buffer([full_checksum[5], full_checksum[13], full_checksum[21], full_checksum[29]]);
	checksum := append(make([]byte, 0, 4), fc[5], fc[13], fc[21], fc[29])
	return checksum
}

func cHash(data []byte, chash_length int) string {
	console.Log("cHash: %#v", data)
	checkLength(chash_length)
//	var hash = crypto.createHash((chash_length === 160) ? "ripemd160" : "sha256").update(data, "utf8").digest();
//	var truncated_hash = (chash_length === 160) ? hash.slice(4) : hash; // drop first 4 bytes if 160
	var hash []byte
	var truncated_hash []byte
	switch chash_length {
	case CHASH160_BITS:
		// [fyi] alas, no Sum() helper for me
		//hash = ripemd160.Sum(data)
		// [tbd] .Reset and reuse
		rmd := ripemd160.New()
		rmd.Write(data)
		hash_ := rmd.Sum(nil)
		hash = hash_[:]
		// drop first 4 bytes if 160
		truncated_hash = hash[4:]
	case CHASH288_BITS:
		hash_ := sha256.Sum256(data)
		hash = hash_[:]
		truncated_hash = hash
	}
	console.Log("hash %#v", hash)
	console.Log("clean data %#v", truncated_hash)
	checksum := checksum(truncated_hash)
	console.Log("checksum %#v", checksum)
	console.Log("checksum bits %#v", buffer2bin(checksum))
	
	binCleanData := buffer2bin(truncated_hash)
	binChecksum := buffer2bin(checksum)
	//console.Log("%d %d", len(binCleanData), len(binChecksum))
	binChash := mixChecksumIntoCleanData(binCleanData, binChecksum)
	//console.log(binCleanData.length, binChecksum.length, binChash.length);
	chash := bin2buffer(binChash)
	console.Log("cHash %#v", chash)
//	var encoded = (chash_length === 160) ? base32.encode(chash).toString() : chash.toString('base64');
	var encoded string
	switch chash_length {
	case CHASH160_BITS:
		encoded = base32.StdEncoding.EncodeToString(chash)
	case CHASH288_BITS:
		encoded = base64.StdEncoding.EncodeToString(chash)
	}
	console.Log("encoded %s", encoded)
	return encoded;
}

func CHash160(data []byte) CHash160T {
	return CHash160T(cHash(data, CHASH160_BITS))
}

func CHash288(data []byte) string {
	return cHash(data, CHASH288_BITS)
}

/**
function isChashValid(encoded){
	var encoded_len = encoded.length;
	if (encoded_len !== 32 && encoded_len !== 48) // 160/5 = 32, 288/6 = 48
		throw Error("wrong encoded length: "+encoded_len);
	try{
		var chash = (encoded_len === 32) ? base32.decode(encoded) : new Buffer(encoded, 'base64');
	}
	catch(e){
		console.log(e);
		return false;
	}
	var binChash = buffer2bin(chash);
	var separated = separateIntoCleanDataAndChecksum(binChash);
	var clean_data = bin2buffer(separated.clean_data);
	//console.log("clean data", clean_data);
	var checksum = bin2buffer(separated.checksum);
	//console.log(checksum);
	//console.log(getChecksum(clean_data));
	return checksum.equals(getChecksum(clean_data));
}
 **/


const Boo = `
/*jslint node: true * /
"use strict";
var crypto = require('crypto');
var base32 = require('thirty-two');

var PI = "14159265358979323846264338327950288419716939937510";
var zeroString = "00000000";

var arrRelativeOffsets = PI.split("");

function checkLength(chash_length){
	if (chash_length !== 160 && chash_length !== 288)
		throw Error("unsupported c-hash length: "+chash_length);
}

function calcOffsets(chash_length){
	checkLength(chash_length);
	var arrOffsets = [];
	var offset = 0;
	var index = 0;

	for (var i=0; offset<chash_length; i++){
		var relative_offset = parseInt(arrRelativeOffsets[i]);
		if (relative_offset === 0)
			continue;
		offset += relative_offset;
		if (chash_length === 288)
			offset += 4;
		if (offset >= chash_length)
			break;
		arrOffsets.push(offset);
		//console.log("index="+index+", offset="+offset);
		index++;
	}

	if (index != 32)
		throw Error("wrong number of checksum bits");
	
	return arrOffsets;
}

var arrOffsets160 = calcOffsets(160);
var arrOffsets288 = calcOffsets(288);

function separateIntoCleanDataAndChecksum(bin){
	var len = bin.length;
	var arrOffsets;
	if (len === 160)
		arrOffsets = arrOffsets160;
	else if (len === 288)
		arrOffsets = arrOffsets288;
	else
		throw Error("bad length="+len+", bin = "+bin);
	var arrFrags = [];
	var arrChecksumBits = [];
	var start = 0;
	for (var i=0; i<arrOffsets.length; i++){
		arrFrags.push(bin.substring(start, arrOffsets[i]));
		arrChecksumBits.push(bin.substr(arrOffsets[i], 1));
		start = arrOffsets[i]+1;
	}
	// add last frag
	if (start < bin.length)
		arrFrags.push(bin.substring(start));
	var binCleanData = arrFrags.join("");
	var binChecksum = arrChecksumBits.join("");
	return {clean_data: binCleanData, checksum: binChecksum};
}

function mixChecksumIntoCleanData(binCleanData, binChecksum){
	if (binChecksum.length !== 32)
		throw Error("bad checksum length");
	var len = binCleanData.length + binChecksum.length;
	var arrOffsets;
	if (len === 160)
		arrOffsets = arrOffsets160;
	else if (len === 288)
		arrOffsets = arrOffsets288;
	else
		throw Error("bad length="+len+", clean data = "+binCleanData+", checksum = "+binChecksum);
	var arrFrags = [];
	var arrChecksumBits = binChecksum.split("");
	var start = 0;
	for (var i=0; i<arrOffsets.length; i++){
		var end = arrOffsets[i] - i;
		arrFrags.push(binCleanData.substring(start, end));
		arrFrags.push(arrChecksumBits[i]);
		start = end;
	}
	// add last frag
	if (start < binCleanData.length)
		arrFrags.push(binCleanData.substring(start));
	return arrFrags.join("");
}

function buffer2bin(buf){
	var bytes = [];
	for (var i=0; i<buf.length; i++){
		var bin = buf[i].toString(2);
		if (bin.length < 8) // pad with zeros
			bin = zeroString.substring(bin.length, 8) + bin;
		bytes.push(bin);
	}
	return bytes.join("");
}

function bin2buffer(bin){
	var len = bin.length/8;
	var buf = new Buffer(len);
	for (var i=0; i<len; i++)
		buf[i] = parseInt(bin.substr(i*8, 8), 2);
	return buf;
}

function getChecksum(clean_data){
	var full_checksum = crypto.createHash("sha256").update(clean_data).digest();
	//console.log(full_checksum);
	var checksum = new Buffer([full_checksum[5], full_checksum[13], full_checksum[21], full_checksum[29]]);
	return checksum;
}

function getChash(data, chash_length){
	//console.log("getChash: "+data);
	checkLength(chash_length);
	var hash = crypto.createHash((chash_length === 160) ? "ripemd160" : "sha256").update(data, "utf8").digest();
	//console.log("hash", hash);
	var truncated_hash = (chash_length === 160) ? hash.slice(4) : hash; // drop first 4 bytes if 160
	//console.log("clean data", truncated_hash);
	var checksum = getChecksum(truncated_hash);
	//console.log("checksum", checksum);
	//console.log("checksum", buffer2bin(checksum));
	
	var binCleanData = buffer2bin(truncated_hash);
	var binChecksum = buffer2bin(checksum);
	var binChash = mixChecksumIntoCleanData(binCleanData, binChecksum);
	//console.log(binCleanData.length, binChecksum.length, binChash.length);
	var chash = bin2buffer(binChash);
	//console.log("chash     ", chash);
	var encoded = (chash_length === 160) ? base32.encode(chash).toString() : chash.toString('base64');
	//console.log(encoded);
	return encoded;
}

function getChash160(data){
	return getChash(data, 160);
}

function getChash288(data){
	return getChash(data, 288);
}

function isChashValid(encoded){
	var encoded_len = encoded.length;
	if (encoded_len !== 32 && encoded_len !== 48) // 160/5 = 32, 288/6 = 48
		throw Error("wrong encoded length: "+encoded_len);
	try{
		var chash = (encoded_len === 32) ? base32.decode(encoded) : new Buffer(encoded, 'base64');
	}
	catch(e){
		console.log(e);
		return false;
	}
	var binChash = buffer2bin(chash);
	var separated = separateIntoCleanDataAndChecksum(binChash);
	var clean_data = bin2buffer(separated.clean_data);
	//console.log("clean data", clean_data);
	var checksum = bin2buffer(separated.checksum);
	//console.log(checksum);
	//console.log(getChecksum(clean_data));
	return checksum.equals(getChecksum(clean_data));
}


exports.getChash160 = getChash160;
exports.getChash288 = getChash288;
exports.isChashValid = isChashValid;

`
