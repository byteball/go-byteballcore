package constants

import(
	"math"

 .	"github.com/byteball/go-byteballcore/types"
)

const(
	COUNT_WITNESSES		int = 12
	MAX_WITNESS_LIST_MUTATIONS int = 1
	TOTAL_WHITEBYTES	float64 = 1e15

	COUNT_MC_BALLS_FOR_PAID_WITNESSING int = 100

	Version			string = "1.0"
	Alt			string = "1"

	HASH_LENGTH		int = 44
	PUBKEY_LENGTH		int = 44
	SIG_LENGTH		int = 88


	MAX_PROFILE_FIELD_LENGTH	int = 50
	MAX_PROFILE_VALUE_LENGTH	int = 100


)

type(
)

var(
	MAJORITY_OF_WITNESSES	int

	GENESIS_UNIT		UnitT
	BLACKBYTES_ASSET	AssetT
)

func init() {
	if COUNT_WITNESSES%2==0 {
		MAJORITY_OF_WITNESSES = (COUNT_WITNESSES/2+1)
	} else {
		MAJORITY_OF_WITNESSES  = int(math.Ceil(float64(COUNT_WITNESSES)/2))
	}

	if Alt == "2" && Version == "1.0t" {
		GENESIS_UNIT = "TvqutGPz3T4Cs6oiChxFlclY92M2MvCvfXR5/FETato="
		BLACKBYTES_ASSET = "LUQu5ik4WLfCrr8OwXezqBa+i3IlZLqxj2itQZQm8WY="
	} else {
		GENESIS_UNIT = "oj8yEksX9Ubq7lLc+p6F2uyHUuynugeVq4+ikT67X6E="
		BLACKBYTES_ASSET = "qO2JsiuDMh/j+pqJYZw3u82O71WjCDf0vTNvsnntr8o="
	}

}

/**
"use strict";

exports.COUNT_WITNESSES = 12;
exports.MAX_WITNESS_LIST_MUTATIONS = 1;
exports.TOTAL_WHITEBYTES = 1e15;
exports.MAJORITY_OF_WITNESSES = (exports.COUNT_WITNESSES%2===0) ? (exports.COUNT_WITNESSES/2+1) : Math.ceil(exports.COUNT_WITNESSES/2);
exports.COUNT_MC_BALLS_FOR_PAID_WITNESSING = 100;

exports.version = '1.0';
exports.alt = '1';

exports.GENESIS_UNIT = (exports.alt === '2' && exports.version === '1.0t') ? 'TvqutGPz3T4Cs6oiChxFlclY92M2MvCvfXR5/FETato=' : 'oj8yEksX9Ubq7lLc+p6F2uyHUuynugeVq4+ikT67X6E=';
exports.BLACKBYTES_ASSET = (exports.alt === '2' && exports.version === '1.0t') ? 'LUQu5ik4WLfCrr8OwXezqBa+i3IlZLqxj2itQZQm8WY=' : 'qO2JsiuDMh/j+pqJYZw3u82O71WjCDf0vTNvsnntr8o=';

exports.HASH_LENGTH = 44;
exports.PUBKEY_LENGTH = 44;
exports.SIG_LENGTH = 88;

// anti-spam limits
exports.MAX_AUTHORS_PER_UNIT = 16;
exports.MAX_PARENTS_PER_UNIT = 16;
exports.MAX_MESSAGES_PER_UNIT = 128;
exports.MAX_SPEND_PROOFS_PER_MESSAGE = 128;
exports.MAX_INPUTS_PER_PAYMENT_MESSAGE = 128;
exports.MAX_OUTPUTS_PER_PAYMENT_MESSAGE = 128;
exports.MAX_CHOICES_PER_POLL = 128;
exports.MAX_DENOMINATIONS_PER_ASSET_DEFINITION = 64;
exports.MAX_ATTESTORS_PER_ASSET = 64;
exports.MAX_DATA_FEED_NAME_LENGTH = 64;
exports.MAX_DATA_FEED_VALUE_LENGTH = 64;
exports.MAX_AUTHENTIFIER_LENGTH = 4096;
exports.MAX_CAP = 9e15;
exports.MAX_COMPLEXITY = 100;

exports.MAX_PROFILE_FIELD_LENGTH = 50;
exports.MAX_PROFILE_VALUE_LENGTH = 100;

exports.TEXTCOIN_CLAIM_FEE = 548;
exports.TEXTCOIN_ASSET_CLAIM_FEE = 750;
exports.TEXTCOIN_ASSET_CLAIM_HEADER_FEE = 391;
exports.TEXTCOIN_ASSET_CLAIM_MESSAGE_FEE = 209;
exports.TEXTCOIN_ASSET_CLAIM_BASE_MSG_FEE = 158;
 **/
