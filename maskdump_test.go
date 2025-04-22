package main

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestMaskEmailWithRules(t *testing.T) {
	typeMaskingInfo := Email

	// Example #1 target and value
	email1 := "test@example.com"
	target1 := "username:2-"
	value1 := "hash:6"

	parts1 := strings.Split(email1, "@")
	if len(parts1) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email1)
	}

	localPart1 := parts1[0]
	domainPart1 := parts1[1]

	positions1 := parseTargetPositions(strings.TrimPrefix(target1, "username:"), len(localPart1))
	maskLocalPart1 := applyMasking(localPart1, positions1, value1, typeMaskingInfo)

	maskedEmail1 := maskLocalPart1 + "@" + domainPart1
	if maskedEmail1 != "t098f6b@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email1, target1, value1, maskedEmail1)
	}

	// Example #2 target and value
	email2 := "test@example.com"
	target2 := "domain:2-"
	value2 := "hash:6"

	parts2 := strings.Split(email2, "@")
	if len(parts2) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email2)
	}

	localPart2 := parts2[0]
	domainPart2 := parts2[1]

	positions2 := parseTargetPositions(strings.TrimPrefix(target2, "domain:"), len(domainPart2))
	maskDomainPart2 := applyMasking(domainPart2, positions2, value2, typeMaskingInfo)

	maskedEmail2 := localPart2 + "@" + maskDomainPart2
	if maskedEmail2 != "test@e5ababd" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email2, target2, value2, maskedEmail2)
	}

	// Example #3 target and value
	email3 := "test@example.com"
	target3 := "2-"
	value3 := "hash:6"

	positions3 := parseTargetPositions(target3, len(email3))
	maskedEmail3 := applyMasking(email3, positions3, value3, typeMaskingInfo)

	if maskedEmail3 != "t55502f" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email3, target3, value3, maskedEmail3)
	}

	// Example #4 target and value
	email4 := "test@example.com"
	target4 := "username:2-"
	value4 := "hash"

	parts4 := strings.Split(email4, "@")
	if len(parts4) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email4)
	}

	localPart4 := parts4[0]
	domainPart4 := parts4[1]

	positions4 := parseTargetPositions(strings.TrimPrefix(target4, "username:"), len(localPart4))
	maskLocalPart4 := applyMasking(localPart4, positions4, value4, typeMaskingInfo)

	maskedEmail4 := maskLocalPart4 + "@" + domainPart4
	if maskedEmail4 != "t098@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email4, target4, value4, maskedEmail4)
	}

	// Example #5 target and value
	email5 := "test@example.com"
	target5 := "domain:2-"
	value5 := "hash"

	parts5 := strings.Split(email5, "@")
	if len(parts5) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email5)
	}

	localPart5 := parts5[0]
	domainPart5 := parts5[1]

	positions5 := parseTargetPositions(strings.TrimPrefix(target5, "domain:"), len(domainPart5))
	maskDomainPart5 := applyMasking(domainPart5, positions5, value5, typeMaskingInfo)

	maskedEmail5 := localPart5 + "@" + maskDomainPart5
	if maskedEmail5 != "test@e5ababd603b" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email5, target5, value5, maskedEmail5)
	}

	// Example #6 target and value
	email6 := "test@example.com"
	target6 := "2-"
	value6 := "hash"

	positions6 := parseTargetPositions(target6, len(email6))
	maskedEmail6 := applyMasking(email6, positions6, value6, typeMaskingInfo)

	if maskedEmail6 != "t55502f40dc8b7c7" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email6, target6, value6, maskedEmail6)
	}

	// Example #7 target and value
	email7 := "test@example.com"
	target7 := "username:2~"
	value7 := "hash:6"

	parts7 := strings.Split(email7, "@")
	if len(parts7) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email7)
	}

	localPart7 := parts7[0]
	domainPart7 := parts7[1]

	positions7 := parseTargetPositions(strings.TrimPrefix(target7, "username:"), len(localPart7))
	maskLocalPart7 := applyMasking(localPart7, positions7, value7, typeMaskingInfo)

	maskedEmail7 := maskLocalPart7 + "@" + domainPart7
	if maskedEmail7 != "te098f6b@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email7, target7, value7, maskedEmail7)
	}

	// Example #8 target and value
	email8 := "test@example.com"
	target8 := "domain:2~"
	value8 := "hash:6"

	parts8 := strings.Split(email8, "@")
	if len(parts8) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email8)
	}

	localPart8 := parts8[0]
	domainPart8 := parts8[1]

	positions8 := parseTargetPositions(strings.TrimPrefix(target8, "domain:"), len(domainPart8))
	maskDomainPart8 := applyMasking(domainPart8, positions8, value8, typeMaskingInfo)

	maskedEmail8 := localPart8 + "@" + maskDomainPart8
	if maskedEmail8 != "test@ex5ababd" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email8, target8, value8, maskedEmail8)
	}

	// Example #9 target and value
	email9 := "test@example.com"
	target9 := "2~"
	value9 := "hash:6"

	positions9 := parseTargetPositions(target9, len(email9))
	maskedEmail9 := applyMasking(email9, positions9, value9, typeMaskingInfo)

	if maskedEmail9 != "te55502f" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email9, target9, value9, maskedEmail9)
	}

	// Example #10 target and value
	email10 := "test@example.com"
	target10 := "username:2~"
	value10 := "hash"

	parts10 := strings.Split(email10, "@")
	if len(parts10) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email10)
	}

	localPart10 := parts10[0]
	domainPart10 := parts10[1]

	positions10 := parseTargetPositions(strings.TrimPrefix(target10, "username:"), len(localPart10))
	maskLocalPart10 := applyMasking(localPart10, positions10, value10, typeMaskingInfo)

	maskedEmail10 := maskLocalPart10 + "@" + domainPart10
	if maskedEmail10 != "te09@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email10, target10, value10, maskedEmail10)
	}

	// Example #11 target and value
	email11 := "test@example.com"
	target11 := "domain:2~"
	value11 := "hash"

	parts11 := strings.Split(email11, "@")
	if len(parts11) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email11)
	}

	localPart11 := parts11[0]
	domainPart11 := parts11[1]

	positions11 := parseTargetPositions(strings.TrimPrefix(target11, "domain:"), len(domainPart11))
	maskDomainPart11 := applyMasking(domainPart11, positions11, value11, typeMaskingInfo)

	maskedEmail11 := localPart11 + "@" + maskDomainPart11
	if maskedEmail11 != "test@ex5ababd603" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email11, target11, value11, maskedEmail11)
	}

	// Example #12 target and value
	email12 := "test@example.com"
	target12 := "2~"
	value12 := "hash"

	positions12 := parseTargetPositions(target12, len(email12))
	maskedEmail12 := applyMasking(email12, positions12, value12, typeMaskingInfo)

	if maskedEmail12 != "te55502f40dc8b7c" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email12, target12, value12, maskedEmail12)
	}

	// Example #13 target and value
	email13 := "test@example.com"
	target13 := "username:-2"
	value13 := "hash:6"

	parts13 := strings.Split(email13, "@")
	if len(parts13) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email13)
	}

	localPart13 := parts13[0]
	domainPart13 := parts13[1]

	positions13 := parseTargetPositions(strings.TrimPrefix(target13, "username:"), len(localPart13))
	maskLocalPart13 := applyMasking(localPart13, positions13, value13, typeMaskingInfo)

	maskedEmail13 := maskLocalPart13 + "@" + domainPart13
	if maskedEmail13 != "098f6bst@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email13, target13, value13, maskedEmail13)
	}

	// Example #14 target and value
	email14 := "test@example.com"
	target14 := "domain:-2"
	value14 := "hash:6"

	parts14 := strings.Split(email14, "@")
	if len(parts14) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email14)
	}

	localPart14 := parts14[0]
	domainPart14 := parts14[1]

	positions14 := parseTargetPositions(strings.TrimPrefix(target14, "domain:"), len(domainPart14))
	maskDomainPart14 := applyMasking(domainPart14, positions14, value14, typeMaskingInfo)

	maskedEmail14 := localPart14 + "@" + maskDomainPart14
	if maskedEmail14 != "test@5ababdample.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email14, target14, value14, maskedEmail14)
	}

	// Example #15 target and value
	email15 := "test@example.com"
	target15 := "-2"
	value15 := "hash:6"

	positions15 := parseTargetPositions(target15, len(email15))
	maskedEmail15 := applyMasking(email15, positions15, value15, typeMaskingInfo)

	if maskedEmail15 != "55502fst@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email15, target15, value15, maskedEmail15)
	}

	// Example #16 target and value
	email16 := "test@example.com"
	target16 := "username:-2"
	value16 := "hash"

	parts16 := strings.Split(email16, "@")
	if len(parts16) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email16)
	}

	localPart16 := parts16[0]
	domainPart16 := parts16[1]

	positions16 := parseTargetPositions(strings.TrimPrefix(target16, "username:"), len(localPart16))
	maskLocalPart16 := applyMasking(localPart16, positions16, value16, typeMaskingInfo)

	maskedEmail16 := maskLocalPart16 + "@" + domainPart16
	if maskedEmail16 != "09st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email16, target16, value16, maskedEmail16)
	}

	// Example #17 target and value
	email17 := "test@example.com"
	target17 := "domain:-2"
	value17 := "hash"

	parts17 := strings.Split(email17, "@")
	if len(parts17) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email17)
	}

	localPart17 := parts17[0]
	domainPart17 := parts17[1]

	positions17 := parseTargetPositions(strings.TrimPrefix(target17, "domain:"), len(domainPart17))
	maskDomainPart17 := applyMasking(domainPart17, positions17, value17, typeMaskingInfo)

	maskedEmail17 := localPart17 + "@" + maskDomainPart17
	if maskedEmail17 != "test@5aample.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email17, target17, value17, maskedEmail17)
	}

	// Example #18 target and value
	email18 := "test@example.com"
	target18 := "-2"
	value18 := "hash"

	positions18 := parseTargetPositions(target18, len(email18))
	maskedEmail18 := applyMasking(email18, positions18, value18, typeMaskingInfo)

	if maskedEmail18 != "55st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email18, target18, value18, maskedEmail18)
	}

	// Example #19 target and value
	email19 := "test@example.com"
	target19 := "username:~2"
	value19 := "hash:6"

	parts19 := strings.Split(email19, "@")
	if len(parts19) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email19)
	}

	localPart19 := parts19[0]
	domainPart19 := parts19[1]

	positions19 := parseTargetPositions(strings.TrimPrefix(target19, "username:"), len(localPart19))
	maskLocalPart19 := applyMasking(localPart19, positions19, value19, typeMaskingInfo)

	maskedEmail19 := maskLocalPart19 + "@" + domainPart19
	if maskedEmail19 != "098f6bst@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email19, target19, value19, maskedEmail19)
	}

	// Example #20 target and value
	email20 := "test@example.com"
	target20 := "domain:~2"
	value20 := "hash:6"

	parts20 := strings.Split(email20, "@")
	if len(parts20) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email20)
	}

	localPart20 := parts20[0]
	domainPart20 := parts20[1]

	positions20 := parseTargetPositions(strings.TrimPrefix(target20, "domain:"), len(domainPart20))
	maskDomainPart20 := applyMasking(domainPart20, positions20, value20, typeMaskingInfo)

	maskedEmail20 := localPart20 + "@" + maskDomainPart20
	if maskedEmail20 != "test@5ababdom" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email20, target20, value20, maskedEmail20)
	}

	// Example #21 target and value
	email21 := "test@example.com"
	target21 := "~2"
	value21 := "hash:6"

	positions21 := parseTargetPositions(target21, len(email21))
	maskedEmail21 := applyMasking(email21, positions21, value21, typeMaskingInfo)

	if maskedEmail21 != "55502fom" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email21, target21, value21, maskedEmail21)
	}

	// Example #22 target and value
	email22 := "test@example.com"
	target22 := "username:~2"
	value22 := "hash"

	parts22 := strings.Split(email22, "@")
	if len(parts22) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email22)
	}

	localPart22 := parts22[0]
	domainPart22 := parts22[1]

	positions22 := parseTargetPositions(strings.TrimPrefix(target22, "username:"), len(localPart22))
	maskLocalPart22 := applyMasking(localPart22, positions22, value22, typeMaskingInfo)

	maskedEmail22 := maskLocalPart22 + "@" + domainPart22
	if maskedEmail22 != "09st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email22, target22, value22, maskedEmail22)
	}

	// Example #23 target and value
	email23 := "test@example.com"
	target23 := "domain:~2"
	value23 := "hash"

	parts23 := strings.Split(email23, "@")
	if len(parts23) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email23)
	}

	localPart23 := parts23[0]
	domainPart23 := parts23[1]

	positions23 := parseTargetPositions(strings.TrimPrefix(target23, "domain:"), len(domainPart23))
	maskDomainPart23 := applyMasking(domainPart23, positions23, value23, typeMaskingInfo)

	maskedEmail23 := localPart23 + "@" + maskDomainPart23
	if maskedEmail23 != "test@5ababd603om" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email23, target23, value23, maskedEmail23)
	}

	// Example #24 target and value
	email24 := "test@example.com"
	target24 := "~2"
	value24 := "hash"

	positions24 := parseTargetPositions(target24, len(email24))
	maskedEmail24 := applyMasking(email24, positions24, value24, typeMaskingInfo)

	if maskedEmail24 != "55502f40dc8b7com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email24, target24, value24, maskedEmail24)
	}

	// Example #25 target and value
	email25 := "test@example.com"
	target25 := "username:2-5"
	value25 := "hash:6"

	parts25 := strings.Split(email25, "@")
	if len(parts25) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email25)
	}

	localPart25 := parts25[0]
	domainPart25 := parts25[1]

	positions25 := parseTargetPositions(strings.TrimPrefix(target25, "username:"), len(localPart25))
	maskLocalPart25 := applyMasking(localPart25, positions25, value25, typeMaskingInfo)

	maskedEmail25 := maskLocalPart25 + "@" + domainPart25
	if maskedEmail25 != "t098f6b@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email25, target25, value25, maskedEmail25)
	}

	// Example #26 target and value
	email26 := "test@example.com"
	target26 := "domain:2-5"
	value26 := "hash:6"

	parts26 := strings.Split(email26, "@")
	if len(parts26) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email26)
	}

	localPart26 := parts26[0]
	domainPart26 := parts26[1]

	positions26 := parseTargetPositions(strings.TrimPrefix(target26, "domain:"), len(domainPart26))
	maskDomainPart26 := applyMasking(domainPart26, positions26, value26, typeMaskingInfo)

	maskedEmail26 := localPart26 + "@" + maskDomainPart26
	if maskedEmail26 != "test@e5ababdle.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email26, target26, value26, maskedEmail26)
	}

	// Example #27 target and value
	email27 := "test@example.com"
	target27 := "2-5"
	value27 := "hash:6"

	positions27 := parseTargetPositions(target27, len(email27))
	maskedEmail27 := applyMasking(email27, positions27, value27, typeMaskingInfo)

	if maskedEmail27 != "t55502fexample.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email27, target27, value27, maskedEmail27)
	}

	// Example #28 target and value
	email28 := "test@example.com"
	target28 := "username:2-5"
	value28 := "hash"

	parts28 := strings.Split(email28, "@")
	if len(parts28) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email28)
	}

	localPart28 := parts28[0]
	domainPart28 := parts28[1]

	positions28 := parseTargetPositions(strings.TrimPrefix(target28, "username:"), len(localPart28))
	maskLocalPart28 := applyMasking(localPart28, positions28, value28, typeMaskingInfo)

	maskedEmail28 := maskLocalPart28 + "@" + domainPart28
	if maskedEmail28 != "t098@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email28, target28, value28, maskedEmail28)
	}

	// Example #29 target and value
	email29 := "test@example.com"
	target29 := "domain:2-5"
	value29 := "hash"

	parts29 := strings.Split(email29, "@")
	if len(parts29) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email29)
	}

	localPart29 := parts29[0]
	domainPart29 := parts29[1]

	positions29 := parseTargetPositions(strings.TrimPrefix(target29, "domain:"), len(domainPart29))
	maskDomainPart29 := applyMasking(domainPart29, positions29, value29, typeMaskingInfo)

	maskedEmail29 := localPart29 + "@" + maskDomainPart29
	if maskedEmail29 != "test@e5abale.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email29, target29, value29, maskedEmail29)
	}

	// Example #30 target and value
	email30 := "test@example.com"
	target30 := "2-5"
	value30 := "hash"

	positions30 := parseTargetPositions(target30, len(email30))
	maskedEmail30 := applyMasking(email30, positions30, value30, typeMaskingInfo)

	if maskedEmail30 != "t5550example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email30, target30, value30, maskedEmail30)
	}

	// Example #31 target and value
	email31 := "test@example.com"
	target31 := "username:2~1"
	value31 := "hash:6"

	parts31 := strings.Split(email31, "@")
	if len(parts31) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email31)
	}

	localPart31 := parts31[0]
	domainPart31 := parts31[1]

	positions31 := parseTargetPositions(strings.TrimPrefix(target31, "username:"), len(localPart31))
	maskLocalPart31 := applyMasking(localPart31, positions31, value31, typeMaskingInfo)

	maskedEmail31 := maskLocalPart31 + "@" + domainPart31
	if maskedEmail31 != "te098f6bt@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email31, target31, value31, maskedEmail31)
	}

	// Example #32 target and value
	email32 := "test@example.com"
	target32 := "domain:2~1"
	value32 := "hash:6"

	parts32 := strings.Split(email32, "@")
	if len(parts32) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email32)
	}

	localPart32 := parts32[0]
	domainPart32 := parts32[1]

	positions32 := parseTargetPositions(strings.TrimPrefix(target32, "domain:"), len(domainPart32))
	maskDomainPart32 := applyMasking(domainPart32, positions32, value32, typeMaskingInfo)

	maskedEmail32 := localPart32 + "@" + maskDomainPart32
	if maskedEmail32 != "test@ex5ababdm" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email32, target32, value32, maskedEmail32)
	}

	// Example #33 target and value
	email33 := "test@example.com"
	target33 := "2~1"
	value33 := "hash:6"

	positions33 := parseTargetPositions(target33, len(email33))
	maskedEmail33 := applyMasking(email33, positions33, value33, typeMaskingInfo)

	if maskedEmail33 != "te55502fm" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email33, target33, value33, maskedEmail33)
	}

	// Example #34 target and value
	email34 := "test@example.com"
	target34 := "username:2~1"
	value34 := "hash"

	parts34 := strings.Split(email34, "@")
	if len(parts34) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email34)
	}

	localPart34 := parts34[0]
	domainPart34 := parts34[1]

	positions34 := parseTargetPositions(strings.TrimPrefix(target34, "username:"), len(localPart34))
	maskLocalPart34 := applyMasking(localPart34, positions34, value34, typeMaskingInfo)

	maskedEmail34 := maskLocalPart34 + "@" + domainPart34
	if maskedEmail34 != "te0t@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email34, target34, value34, maskedEmail34)
	}

	// Example #35 target and value
	email35 := "test@example.com"
	target35 := "domain:2~1"
	value35 := "hash"

	parts35 := strings.Split(email35, "@")
	if len(parts35) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email35)
	}

	localPart35 := parts35[0]
	domainPart35 := parts35[1]

	positions35 := parseTargetPositions(strings.TrimPrefix(target35, "domain:"), len(domainPart35))
	maskDomainPart35 := applyMasking(domainPart35, positions35, value35, typeMaskingInfo)

	maskedEmail35 := localPart35 + "@" + maskDomainPart35
	if maskedEmail35 != "test@ex5ababd60m" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email35, target35, value35, maskedEmail35)
	}

	// Example #36 target and value
	email36 := "test@example.com"
	target36 := "2~1"
	value36 := "hash"

	positions36 := parseTargetPositions(target36, len(email36))
	maskedEmail36 := applyMasking(email36, positions36, value36, typeMaskingInfo)

	if maskedEmail36 != "te55502f40dc8b7m" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email36, target36, value36, maskedEmail36)
	}

	// Tests for symbol *

	// Example #37 target and value
	email37 := "test@example.com"
	target37 := "username:2-"
	value37 := "*"

	parts37 := strings.Split(email37, "@")
	if len(parts37) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email37)
	}

	localPart37 := parts37[0]
	domainPart37 := parts37[1]

	positions37 := parseTargetPositions(strings.TrimPrefix(target37, "username:"), len(localPart37))
	maskLocalPart37 := applyMasking(localPart37, positions37, value37, typeMaskingInfo)

	maskedEmail37 := maskLocalPart37 + "@" + domainPart37
	if maskedEmail37 != "t***@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email37, target37, value37, maskedEmail37)
	}

	// Example #38 target and value
	email38 := "test@example.com"
	target38 := "domain:2-"
	value38 := "*"

	parts38 := strings.Split(email38, "@")
	if len(parts38) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email38)
	}

	localPart38 := parts38[0]
	domainPart38 := parts38[1]

	positions38 := parseTargetPositions(strings.TrimPrefix(target38, "domain:"), len(domainPart38))
	maskDomainPart38 := applyMasking(domainPart38, positions38, value38, typeMaskingInfo)

	maskedEmail38 := localPart38 + "@" + maskDomainPart38
	if maskedEmail38 != "test@e**********" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email38, target38, value38, maskedEmail38)
	}

	// Example #3 target and value
	email39 := "test@example.com"
	target39 := "2-"
	value39 := "*"

	positions39 := parseTargetPositions(target39, len(email39))
	maskedEmail39 := applyMasking(email39, positions39, value39, typeMaskingInfo)

	if maskedEmail39 != "t***************" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email39, target39, value39, maskedEmail39)
	}

	// Example #40 target and value
	email40 := "test@example.com"
	target40 := "username:2~"
	value40 := "*"

	parts40 := strings.Split(email40, "@")
	if len(parts40) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email40)
	}

	localPart40 := parts40[0]
	domainPart40 := parts40[1]

	positions40 := parseTargetPositions(strings.TrimPrefix(target40, "username:"), len(localPart40))
	maskLocalPart40 := applyMasking(localPart40, positions40, value40, typeMaskingInfo)

	maskedEmail40 := maskLocalPart40 + "@" + domainPart40
	if maskedEmail40 != "te**@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email40, target40, value40, maskedEmail40)
	}

	// Example #41 target and value
	email41 := "test@example.com"
	target41 := "domain:2~"
	value41 := "*"

	parts41 := strings.Split(email41, "@")
	if len(parts41) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email41)
	}

	localPart41 := parts41[0]
	domainPart41 := parts41[1]

	positions41 := parseTargetPositions(strings.TrimPrefix(target41, "domain:"), len(domainPart41))
	maskDomainPart41 := applyMasking(domainPart41, positions41, value41, typeMaskingInfo)

	maskedEmail41 := localPart41 + "@" + maskDomainPart41
	if maskedEmail41 != "test@ex*********" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email41, target41, value41, maskedEmail41)
	}

	// Example #42 target and value
	email42 := "test@example.com"
	target42 := "2~"
	value42 := "*"

	positions42 := parseTargetPositions(target42, len(email42))
	maskedEmail42 := applyMasking(email42, positions42, value42, typeMaskingInfo)

	if maskedEmail42 != "te**************" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email42, target42, value42, maskedEmail42)
	}

	// Example #43 target and value
	email43 := "test@example.com"
	target43 := "username:-2"
	value43 := "*"

	parts43 := strings.Split(email43, "@")
	if len(parts43) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email43)
	}

	localPart43 := parts43[0]
	domainPart43 := parts43[1]

	positions43 := parseTargetPositions(strings.TrimPrefix(target43, "username:"), len(localPart43))
	maskLocalPart43 := applyMasking(localPart43, positions43, value43, typeMaskingInfo)

	maskedEmail43 := maskLocalPart43 + "@" + domainPart43
	if maskedEmail43 != "**st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email43, target43, value43, maskedEmail43)
	}

	// Example #44 target and value
	email44 := "test@example.com"
	target44 := "domain:-2"
	value44 := "*"

	parts44 := strings.Split(email44, "@")
	if len(parts44) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email44)
	}

	localPart44 := parts44[0]
	domainPart44 := parts44[1]

	positions44 := parseTargetPositions(strings.TrimPrefix(target44, "domain:"), len(domainPart44))
	maskDomainPart44 := applyMasking(domainPart44, positions44, value44, typeMaskingInfo)

	maskedEmail44 := localPart44 + "@" + maskDomainPart44
	if maskedEmail44 != "test@**ample.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email44, target44, value44, maskedEmail44)
	}

	// Example #45 target and value
	email45 := "test@example.com"
	target45 := "-2"
	value45 := "*"

	positions45 := parseTargetPositions(target45, len(email45))
	maskedEmail45 := applyMasking(email45, positions45, value45, typeMaskingInfo)

	if maskedEmail45 != "**st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email45, target45, value45, maskedEmail45)
	}

	// Example #46 target and value
	email46 := "test@example.com"
	target46 := "username:~2"
	value46 := "*"

	parts46 := strings.Split(email46, "@")
	if len(parts46) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email46)
	}

	localPart46 := parts46[0]
	domainPart46 := parts46[1]

	positions46 := parseTargetPositions(strings.TrimPrefix(target46, "username:"), len(localPart46))
	maskLocalPart46 := applyMasking(localPart46, positions46, value46, typeMaskingInfo)

	maskedEmail46 := maskLocalPart46 + "@" + domainPart46
	if maskedEmail46 != "**st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email46, target46, value46, maskedEmail46)
	}

	// Example #47 target and value
	email47 := "test@example.com"
	target47 := "domain:~2"
	value47 := "*"

	parts47 := strings.Split(email47, "@")
	if len(parts47) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email47)
	}

	localPart47 := parts47[0]
	domainPart47 := parts47[1]

	positions47 := parseTargetPositions(strings.TrimPrefix(target47, "domain:"), len(domainPart47))
	maskDomainPart47 := applyMasking(domainPart47, positions47, value47, typeMaskingInfo)

	maskedEmail47 := localPart47 + "@" + maskDomainPart47
	if maskedEmail47 != "test@*********om" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email47, target47, value47, maskedEmail47)
	}

	// Example #48 target and value
	email48 := "test@example.com"
	target48 := "~2"
	value48 := "*"

	positions48 := parseTargetPositions(target48, len(email48))
	maskedEmail48 := applyMasking(email48, positions48, value48, typeMaskingInfo)

	if maskedEmail48 != "**************om" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email48, target48, value48, maskedEmail48)
	}

	// Example #49 target and value
	email49 := "test@example.com"
	target49 := "username:2-5"
	value49 := "*"

	parts49 := strings.Split(email49, "@")
	if len(parts49) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email49)
	}

	localPart49 := parts49[0]
	domainPart49 := parts49[1]

	positions49 := parseTargetPositions(strings.TrimPrefix(target49, "username:"), len(localPart49))
	maskLocalPart49 := applyMasking(localPart49, positions49, value49, typeMaskingInfo)

	maskedEmail49 := maskLocalPart49 + "@" + domainPart49
	if maskedEmail49 != "t***@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email49, target49, value49, maskedEmail49)
	}

	// Example #50 target and value
	email50 := "test@example.com"
	target50 := "domain:2-5"
	value50 := "*"

	parts50 := strings.Split(email50, "@")
	if len(parts50) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email50)
	}

	localPart50 := parts50[0]
	domainPart50 := parts50[1]

	positions50 := parseTargetPositions(strings.TrimPrefix(target50, "domain:"), len(domainPart50))
	maskDomainPart50 := applyMasking(domainPart50, positions50, value50, typeMaskingInfo)

	maskedEmail50 := localPart50 + "@" + maskDomainPart50
	if maskedEmail50 != "test@e****le.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email50, target50, value50, maskedEmail50)
	}

	// Example #51 target and value
	email51 := "test@example.com"
	target51 := "2-5"
	value51 := "*"

	positions51 := parseTargetPositions(target51, len(email51))
	maskedEmail51 := applyMasking(email51, positions51, value51, typeMaskingInfo)

	if maskedEmail51 != "t****example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email51, target51, value51, maskedEmail51)
	}

	// Example #52 target and value
	email52 := "test@example.com"
	target52 := "username:2-5"
	value52 := "*"

	parts52 := strings.Split(email52, "@")
	if len(parts52) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email52)
	}

	localPart52 := parts52[0]
	domainPart52 := parts52[1]

	positions52 := parseTargetPositions(strings.TrimPrefix(target52, "username:"), len(localPart52))
	maskLocalPart52 := applyMasking(localPart52, positions52, value52, typeMaskingInfo)

	maskedEmail52 := maskLocalPart52 + "@" + domainPart52
	if maskedEmail52 != "t***@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email52, target52, value52, maskedEmail52)
	}

	// Example #53 target and value
	email53 := "test@example.com"
	target53 := "username:2~1"
	value53 := "*"

	parts53 := strings.Split(email53, "@")
	if len(parts53) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email53)
	}

	localPart53 := parts53[0]
	domainPart53 := parts53[1]

	positions53 := parseTargetPositions(strings.TrimPrefix(target53, "username:"), len(localPart53))
	maskLocalPart53 := applyMasking(localPart53, positions53, value53, typeMaskingInfo)

	maskedEmail53 := maskLocalPart53 + "@" + domainPart53
	if maskedEmail53 != "te*t@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email53, target53, value53, maskedEmail53)
	}

	// Example #54 target and value
	email54 := "test@example.com"
	target54 := "domain:2~1"
	value54 := "*"

	parts54 := strings.Split(email54, "@")
	if len(parts54) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email54)
	}

	localPart54 := parts54[0]
	domainPart54 := parts54[1]

	positions54 := parseTargetPositions(strings.TrimPrefix(target54, "domain:"), len(domainPart54))
	maskDomainPart54 := applyMasking(domainPart54, positions54, value54, typeMaskingInfo)

	maskedEmail54 := localPart54 + "@" + maskDomainPart54
	if maskedEmail54 != "test@ex********m" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email54, target54, value54, maskedEmail54)
	}

	// Example #55 target and value
	email55 := "test@example.com"
	target55 := "2~1"
	value55 := "*"

	positions55 := parseTargetPositions(target55, len(email55))
	maskedEmail55 := applyMasking(email55, positions55, value55, typeMaskingInfo)

	if maskedEmail55 != "te*************m" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email55, target55, value55, maskedEmail55)
	}
}

func TestMaskPhoneWithRules(t *testing.T) {
	testCases := []struct {
		phone    string
		target   string
		value    string
		expected string
	}{
		// +7 (900) 111-22-33
		{phone: "+7 (900) 111-22-33", target: "2-6", value: "hash", expected: "+7 (916) 571-22-33"},
		{phone: "+7 (900) 111-22-33", target: "2-6", value: "*", expected: "+7 (***) **1-22-33"},
		{phone: "+7 (900) 111-22-33", target: "2-", value: "hash", expected: "+7 (916) 572-63-35"},
		{phone: "+7 (900) 111-22-33", target: "2-", value: "*", expected: "+7 (***) ***-**-**"},
		{phone: "+7 (900) 111-22-33", target: "-6", value: "hash", expected: "+9 (165) 721-22-33"},
		{phone: "+7 (900) 111-22-33", target: "-6", value: "*", expected: "+* (***) **1-22-33"},
		{phone: "+7 (900) 111-22-33", target: "2~2", value: "hash", expected: "+7 (991) 657-26-33"},
		{phone: "+7 (900) 111-22-33", target: "2~2", value: "*", expected: "+7 (9**) ***-**-33"},
		{phone: "+7 (900) 111-22-33", target: "2~", value: "hash", expected: "+7 (991) 657-26-33"},
		{phone: "+7 (900) 111-22-33", target: "2~", value: "*", expected: "+7 (9**) ***-**-**"},
		{phone: "+7 (900) 111-22-33", target: "~3", value: "hash", expected: "+9 (165) 726-32-33"},
		{phone: "+7 (900) 111-22-33", target: "~3", value: "*", expected: "+* (***) ***-*2-33"},
		{phone: "+7 (900) 111-22-33", target: "1,2,3", value: "hash", expected: "+9 (160) 111-22-33"},
		{phone: "+7 (900) 111-22-33", target: "1,2,3", value: "*", expected: "+* (**0) 111-22-33"},
		{phone: "+7 (900) 111-22-33", target: "1,3,5,7,9", value: "hash", expected: "+9 (910) 615-27-33"},
		{phone: "+7 (900) 111-22-33", target: "1,3,5,7,9", value: "*", expected: "+* (9*0) *1*-2*-33"},

		// +79001112233
		{phone: "+79001112233", target: "2-6", value: "hash", expected: "+79165712233"},
		{phone: "+79001112233", target: "2-6", value: "*", expected: "+7*****12233"},
		{phone: "+79001112233", target: "2-", value: "hash", expected: "+79165726335"},
		{phone: "+79001112233", target: "2-", value: "*", expected: "+7**********"},
		{phone: "+79001112233", target: "-6", value: "hash", expected: "+91657212233"},
		{phone: "+79001112233", target: "-6", value: "*", expected: "+******12233"},
		{phone: "+79001112233", target: "2~2", value: "hash", expected: "+79916572633"},
		{phone: "+79001112233", target: "2~2", value: "*", expected: "+79*******33"},
		{phone: "+79001112233", target: "2~", value: "hash", expected: "+79916572633"},
		{phone: "+79001112233", target: "2~", value: "*", expected: "+79*********"},
		{phone: "+79001112233", target: "~3", value: "hash", expected: "+91657263233"},
		{phone: "+79001112233", target: "~3", value: "*", expected: "+********233"},
		{phone: "+79001112233", target: "1,2,3", value: "hash", expected: "+91601112233"},
		{phone: "+79001112233", target: "1,2,3", value: "*", expected: "+***01112233"},
		{phone: "+79001112233", target: "1,3,5,7,9", value: "hash", expected: "+99106152733"},
		{phone: "+79001112233", target: "1,3,5,7,9", value: "*", expected: "+*9*0*1*2*33"},

		// 79001112233
		{phone: "79001112233", target: "2-6", value: "hash", expected: "79165712233"},
		{phone: "79001112233", target: "2-6", value: "*", expected: "7*****12233"},
		{phone: "79001112233", target: "2-", value: "hash", expected: "79165726335"},
		{phone: "79001112233", target: "2-", value: "*", expected: "7**********"},
		{phone: "79001112233", target: "-6", value: "hash", expected: "91657212233"},
		{phone: "79001112233", target: "-6", value: "*", expected: "******12233"},
		{phone: "79001112233", target: "2~2", value: "hash", expected: "79916572633"},
		{phone: "79001112233", target: "2~2", value: "*", expected: "79*******33"},
		{phone: "79001112233", target: "2~", value: "hash", expected: "79916572633"},
		{phone: "79001112233", target: "2~", value: "*", expected: "79*********"},
		{phone: "79001112233", target: "~3", value: "hash", expected: "91657263233"},
		{phone: "79001112233", target: "~3", value: "*", expected: "********233"},
		{phone: "79001112233", target: "1,2,3", value: "hash", expected: "91601112233"},
		{phone: "79001112233", target: "1,2,3", value: "*", expected: "***01112233"},
		{phone: "79001112233", target: "1,3,5,7,9", value: "hash", expected: "99106152733"},
		{phone: "79001112233", target: "1,3,5,7,9", value: "*", expected: "*9*0*1*2*33"},

		// 89001112233
		{phone: "89001112233", target: "2-6", value: "hash", expected: "89242412233"},
		{phone: "89001112233", target: "2-6", value: "*", expected: "8*****12233"},
		{phone: "89001112233", target: "2-", value: "hash", expected: "89242485674"},
		{phone: "89001112233", target: "2-", value: "*", expected: "8**********"},
		{phone: "89001112233", target: "-6", value: "hash", expected: "92424812233"},
		{phone: "89001112233", target: "-6", value: "*", expected: "******12233"},
		{phone: "89001112233", target: "2~2", value: "hash", expected: "89924248533"},
		{phone: "89001112233", target: "2~2", value: "*", expected: "89*******33"},
		{phone: "89001112233", target: "2~", value: "hash", expected: "89924248567"},
		{phone: "89001112233", target: "2~", value: "*", expected: "89*********"},
		{phone: "89001112233", target: "~3", value: "hash", expected: "92424856233"},
		{phone: "89001112233", target: "~3", value: "*", expected: "********233"},
		{phone: "89001112233", target: "1,2,3", value: "hash", expected: "92401112233"},
		{phone: "89001112233", target: "1,2,3", value: "*", expected: "***01112233"},
		{phone: "89001112233", target: "1,3,5,7,9", value: "hash", expected: "99204122433"},
		{phone: "89001112233", target: "1,3,5,7,9", value: "*", expected: "*9*0*1*2*33"},

		// 8 (900) 111-22-33
		{phone: "8 (900) 111-22-33", target: "2-6", value: "hash", expected: "8 (924) 241-22-33"},
		{phone: "8 (900) 111-22-33", target: "2-6", value: "*", expected: "8 (***) **1-22-33"},
		{phone: "8 (900) 111-22-33", target: "2-", value: "hash", expected: "8 (924) 248-56-74"},
		{phone: "8 (900) 111-22-33", target: "2-", value: "*", expected: "8 (***) ***-**-**"},
		{phone: "8 (900) 111-22-33", target: "-6", value: "hash", expected: "9 (242) 481-22-33"},
		{phone: "8 (900) 111-22-33", target: "-6", value: "*", expected: "* (***) **1-22-33"},
		{phone: "8 (900) 111-22-33", target: "2~2", value: "hash", expected: "8 (992) 424-85-33"},
		{phone: "8 (900) 111-22-33", target: "2~2", value: "*", expected: "8 (9**) ***-**-33"},
		{phone: "8 (900) 111-22-33", target: "2~", value: "hash", expected: "8 (992) 424-85-67"},
		{phone: "8 (900) 111-22-33", target: "2~", value: "*", expected: "8 (9**) ***-**-**"},
		{phone: "8 (900) 111-22-33", target: "~3", value: "hash", expected: "9 (242) 485-62-33"},
		{phone: "8 (900) 111-22-33", target: "~3", value: "*", expected: "* (***) ***-*2-33"},
		{phone: "8 (900) 111-22-33", target: "1,2,3", value: "hash", expected: "9 (240) 111-22-33"},
		{phone: "8 (900) 111-22-33", target: "1,2,3", value: "*", expected: "* (**0) 111-22-33"},
		{phone: "8 (900) 111-22-33", target: "1,3,5,7,9", value: "hash", expected: "9 (920) 412-24-33"},
		{phone: "8 (900) 111-22-33", target: "1,3,5,7,9", value: "*", expected: "* (9*0) *1*-2*-33"},

		// 8-900-111-22-33
		{phone: "8-900-111-22-33", target: "2-6", value: "hash", expected: "8-924-241-22-33"},
		{phone: "8-900-111-22-33", target: "2-6", value: "*", expected: "8-***-**1-22-33"},
		{phone: "8-900-111-22-33", target: "2-", value: "hash", expected: "8-924-248-56-74"},
		{phone: "8-900-111-22-33", target: "2-", value: "*", expected: "8-***-***-**-**"},
		{phone: "8-900-111-22-33", target: "-6", value: "hash", expected: "9-242-481-22-33"},
		{phone: "8-900-111-22-33", target: "-6", value: "*", expected: "*-***-**1-22-33"},
		{phone: "8-900-111-22-33", target: "2~2", value: "hash", expected: "8-992-424-85-33"},
		{phone: "8-900-111-22-33", target: "2~2", value: "*", expected: "8-9**-***-**-33"},
		{phone: "8-900-111-22-33", target: "2~", value: "hash", expected: "8-992-424-85-67"},
		{phone: "8-900-111-22-33", target: "2~", value: "*", expected: "8-9**-***-**-**"},
		{phone: "8-900-111-22-33", target: "~3", value: "hash", expected: "9-242-485-62-33"},
		{phone: "8-900-111-22-33", target: "~3", value: "*", expected: "*-***-***-*2-33"},
		{phone: "8-900-111-22-33", target: "1,2,3", value: "hash", expected: "9-240-111-22-33"},
		{phone: "8-900-111-22-33", target: "1,2,3", value: "*", expected: "*-**0-111-22-33"},
		{phone: "8-900-111-22-33", target: "1,3,5,7,9", value: "hash", expected: "9-920-412-24-33"},
		{phone: "8-900-111-22-33", target: "1,3,5,7,9", value: "*", expected: "*-9*0-*1*-2*-33"},

		// +7(900)111-22-33
		{phone: "+7(900)111-22-33", target: "2-6", value: "hash", expected: "+7(916)571-22-33"},
		{phone: "+7(900)111-22-33", target: "2-6", value: "*", expected: "+7(***)**1-22-33"},
		{phone: "+7(900)111-22-33", target: "2-", value: "hash", expected: "+7(916)572-63-35"},
		{phone: "+7(900)111-22-33", target: "2-", value: "*", expected: "+7(***)***-**-**"},
		{phone: "+7(900)111-22-33", target: "-6", value: "hash", expected: "+9(165)721-22-33"},
		{phone: "+7(900)111-22-33", target: "-6", value: "*", expected: "+*(***)**1-22-33"},
		{phone: "+7(900)111-22-33", target: "2~2", value: "hash", expected: "+7(991)657-26-33"},
		{phone: "+7(900)111-22-33", target: "2~2", value: "*", expected: "+7(9**)***-**-33"},
		{phone: "+7(900)111-22-33", target: "2~", value: "hash", expected: "+7(991)657-26-33"},
		{phone: "+7(900)111-22-33", target: "2~", value: "*", expected: "+7(9**)***-**-**"},
		{phone: "+7(900)111-22-33", target: "~3", value: "hash", expected: "+9(165)726-32-33"},
		{phone: "+7(900)111-22-33", target: "~3", value: "*", expected: "+*(***)***-*2-33"},
		{phone: "+7(900)111-22-33", target: "1,2,3", value: "hash", expected: "+9(160)111-22-33"},
		{phone: "+7(900)111-22-33", target: "1,2,3", value: "*", expected: "+*(**0)111-22-33"},
		{phone: "+7(900)111-22-33", target: "1,3,5,7,9", value: "hash", expected: "+9(910)615-27-33"},
		{phone: "+7(900)111-22-33", target: "1,3,5,7,9", value: "*", expected: "+*(9*0)*1*-2*-33"},

		// 7(900)111-22-33
		{phone: "7(900)111-22-33", target: "2-6", value: "hash", expected: "7(916)571-22-33"},
		{phone: "7(900)111-22-33", target: "2-6", value: "*", expected: "7(***)**1-22-33"},
		{phone: "7(900)111-22-33", target: "2-", value: "hash", expected: "7(916)572-63-35"},
		{phone: "7(900)111-22-33", target: "2-", value: "*", expected: "7(***)***-**-**"},
		{phone: "7(900)111-22-33", target: "-6", value: "hash", expected: "9(165)721-22-33"},
		{phone: "7(900)111-22-33", target: "-6", value: "*", expected: "*(***)**1-22-33"},
		{phone: "7(900)111-22-33", target: "2~2", value: "hash", expected: "7(991)657-26-33"},
		{phone: "7(900)111-22-33", target: "2~2", value: "*", expected: "7(9**)***-**-33"},
		{phone: "7(900)111-22-33", target: "2~", value: "hash", expected: "7(991)657-26-33"},
		{phone: "7(900)111-22-33", target: "2~", value: "*", expected: "7(9**)***-**-**"},
		{phone: "7(900)111-22-33", target: "~3", value: "hash", expected: "9(165)726-32-33"},
		{phone: "7(900)111-22-33", target: "~3", value: "*", expected: "*(***)***-*2-33"},
		{phone: "7(900)111-22-33", target: "1,2,3", value: "hash", expected: "9(160)111-22-33"},
		{phone: "7(900)111-22-33", target: "1,2,3", value: "*", expected: "*(**0)111-22-33"},
		{phone: "7(900)111-22-33", target: "1,3,5,7,9", value: "hash", expected: "9(910)615-27-33"},
		{phone: "7(900)111-22-33", target: "1,3,5,7,9", value: "*", expected: "*(9*0)*1*-2*-33"},

		// 8(900)111-22-33
		{phone: "8(900)111-22-33", target: "2-6", value: "hash", expected: "8(924)241-22-33"},
		{phone: "8(900)111-22-33", target: "2-6", value: "*", expected: "8(***)**1-22-33"},
		{phone: "8(900)111-22-33", target: "2-", value: "hash", expected: "8(924)248-56-74"},
		{phone: "8(900)111-22-33", target: "2-", value: "*", expected: "8(***)***-**-**"},
		{phone: "8(900)111-22-33", target: "-6", value: "hash", expected: "9(242)481-22-33"},
		{phone: "8(900)111-22-33", target: "-6", value: "*", expected: "*(***)**1-22-33"},
		{phone: "8(900)111-22-33", target: "2~2", value: "hash", expected: "8(992)424-85-33"},
		{phone: "8(900)111-22-33", target: "2~2", value: "*", expected: "8(9**)***-**-33"},
		{phone: "8(900)111-22-33", target: "2~", value: "hash", expected: "8(992)424-85-67"},
		{phone: "8(900)111-22-33", target: "2~", value: "*", expected: "8(9**)***-**-**"},
		{phone: "8(900)111-22-33", target: "~3", value: "hash", expected: "9(242)485-62-33"},
		{phone: "8(900)111-22-33", target: "~3", value: "*", expected: "*(***)***-*2-33"},
		{phone: "8(900)111-22-33", target: "1,2,3", value: "hash", expected: "9(240)111-22-33"},
		{phone: "8(900)111-22-33", target: "1,2,3", value: "*", expected: "*(**0)111-22-33"},
		{phone: "8(900)111-22-33", target: "1,3,5,7,9", value: "hash", expected: "9(920)412-24-33"},
		{phone: "8(900)111-22-33", target: "1,3,5,7,9", value: "*", expected: "*(9*0)*1*-2*-33"},
	}

	for _, tc := range testCases {
		t.Run(tc.phone+"_"+tc.target+"_"+tc.value, func(t *testing.T) {
			digits := regexp.MustCompile(`\d`).FindAllString(tc.phone, -1)
			digitStr := strings.Join(digits, "")

			positions := parseTargetPositions(tc.target, len(digitStr))
			maskedDigits := applyMasking(digitStr, positions, tc.value, Phone)

			var result strings.Builder
			digitIndex := 0
			for _, c := range tc.phone {
				if c >= '0' && c <= '9' {
					if digitIndex < len(maskedDigits) {
						result.WriteByte(maskedDigits[digitIndex])
						digitIndex++
					}
				} else {
					result.WriteRune(c)
				}
			}

			masked := result.String()
			if masked != tc.expected {
				t.Errorf("Input: phone=%s, target=%s, value=%s:\nExpected: %s\nGet: %s",
					tc.phone, tc.target, tc.value, tc.expected, masked)
			}
		})
	}
}

func TestValidateAlgorithms(t *testing.T) {
	config1 := MaskConfig{
		emailAlgorithm: "light-hash",
		phoneAlgorithm: "light-mask",
	}
	err1 := validateAlgorithms(config1)
	if err1 != nil {
		t.Errorf("validateAlgorithms failed for valid config: %v", err1)
	}

	config2 := MaskConfig{
		emailAlgorithm: "invalid",
		phoneAlgorithm: "light-mask",
	}
	err2 := validateAlgorithms(config2)
	if err2 == nil {
		t.Errorf("validateAlgorithms should have failed for invalid email algorithm")
	}

	config3 := MaskConfig{
		emailAlgorithm: "light-hash",
		phoneAlgorithm: "invalid",
	}
	err3 := validateAlgorithms(config3)
	if err3 == nil {
		t.Errorf("validateAlgorithms should have failed for invalid phone algorithm")
	}
}

func TestLoadSaveCache(t *testing.T) {
	// Setup
	AppConfig.CachePath = "test_cache.json"
	defer os.Remove(AppConfig.CachePath) // Clean up after test

	cache := &Cache{
		Emails: map[string]string{"test@example.com": "masked@example.com"},
		Phones: map[string]string{"123-456-7890": "masked"},
	}

	// Save cache
	err := saveCache(cache)
	if err != nil {
		t.Fatalf("Error saving cache: %v", err)
	}

	// Load cache
	loadedCache, err := loadCache()
	if err != nil {
		t.Fatalf("Error loading cache: %v", err)
	}

	// Compare
	if !reflect.DeepEqual(cache.Emails, loadedCache.Emails) || !reflect.DeepEqual(cache.Phones, loadedCache.Phones) {
		t.Errorf("Cache data mismatch. Expected: %v, Got: %v", cache, loadedCache)
	}
}
