package contracts_test

import (
	"encoding/base64"
	"testing"

	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"github.com/synycboom/tinyman-go-sdk/v1/contracts"
)

var (
	appIDV            = constants.TestnetValidatorAppId
	asset1ID          = uint64(0)
	asset2ID          = uint64(21582668)
	poolLogicAsBase64 = "BCAIAQAAzKalCgMEBQYlJA1EMQkyAxJEMRUyAxJEMSAyAxJEMgQiDUQzAQAxABJEMwEQIQ" +
		"cSRDMBGIGs194dEkQzARkiEjMBGyEEEhA3ARoAgAlib290c3RyYXASEEAAXDMBGSMSRDMBG4ECEjcBGgCABHN3YXASEEACOzMBGyISRDcB" +
		"GgCABG1pbnQSQAE7NwEaAIAEYnVybhJAAZg3ARoAgAZyZWRlZW0SQAJbNwEaAIAEZmVlcxJAAnkAIQYhBSQjEk0yBBJENwEaARclEjcBGg" +
		"IXJBIQRDMCADEAEkQzAhAhBBJEMwIhIxJEMwIiIxwSRDMCIyEHEkQzAiQjEkQzAiWACFRNUE9PTDExEkQzAiZRAA+AD1RpbnltYW5Qb29s" +
		"MS4xIBJEMwIngBNodHRwczovL3RpbnltYW4ub3JnEkQzAikyAxJEMwIqMgMSRDMCKzIDEkQzAiwyAxJEMwMAMQASRDMDECEFEkQzAxElEk" +
		"QzAxQxABJEMwMSIxJEJCMTQAAQMwEBMwIBCDMDAQg1AUIBsTMEADEAEkQzBBAhBRJEMwQRJBJEMwQUMQASRDMEEiMSRDMBATMCAQgzAwEI" +
		"MwQBCDUBQgF8MgQhBhJENwEcATEAE0Q3ARwBMwQUEkQzAgAxABNEMwIUMQASRDMDADMCABJEMwIRJRJEMwMUMwMHMwMQIhJNMQASRDMDES" +
		"MzAxAiEk0kEkQzBAAxABJEMwQUMwIAEkQzAQEzBAEINQFCAREyBCEGEkQ3ARwBMQATRDcBHAEzAhQSRDMDFDMDBzMDECISTTcBHAESRDMC" +
		"ADEAEkQzAhQzBAASRDMCESUSRDMDADEAEkQzAxQzAwczAxAiEk0zBAASRDMDESMzAxAiEk0kEkQzBAAxABNEMwQUMQASRDMBATMCAQgzAw" +
		"EINQFCAJAyBCEFEkQ3ARwBMQATRDMCADcBHAESRDMCADEAE0QzAwAxABJEMwIUMwIHMwIQIhJNMQASRDMDFDMDBzMDECISTTMCABJEMwEB" +
		"MwMBCDUBQgA+MgQhBBJENwEcATEAE0QzAhQzAgczAhAiEk03ARwBEkQzAQEzAgEINQFCABIyBCEEEkQzAQEzAgEINQFCAAAzAAAxABNEMw" +
		"AHMQASRDMACDQBD0M="
)

func TestGetPoolLogicV1_1(t *testing.T) {
	acc, err := contracts.PoolLogicSigAccount(appIDV, asset1ID, asset2ID)
	if err != nil {
		t.Errorf("Unexpected err %s", err.Error())

		return
	}

	poolLogic := base64.StdEncoding.EncodeToString(acc.Lsig.Logic)
	if poolLogic != poolLogicAsBase64 {
		t.Errorf("PoolLogicSigAccount returned wrong logic")
	}
}

func TestGetPoolLogicV1_1Reversed(t *testing.T) {
	acc, err := contracts.PoolLogicSigAccount(appIDV, asset2ID, asset1ID)
	if err != nil {
		t.Errorf("Unexpected err %s", err.Error())

		return
	}

	poolLogic := base64.StdEncoding.EncodeToString(acc.Lsig.Logic)
	if poolLogic != poolLogicAsBase64 {
		t.Errorf("PoolLogicSigAccount returned wrong logic")
	}
}
