package googleplay

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestVerify(t *testing.T) {
	googlePlay, err := os.ReadFile("./sprintsguys-d0a7a27c99d0.json")
	if err != nil {
		panic(err)
	}
	client, err := New(googlePlay)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(client.VerifyProduct(context.Background(), "com.lqyy.royal.sprintsguys.android", "1010004", "imilffoffcccieknfbmnjgok.AO-J1OxOWu7NKu3cNdmJghyPRKle3Nxq2k9qp-gRBLf1CpdWN3KTvig6GwPdYUtB4Ey4wSA8QYYwgV0NWJvdFfvEm7prDIGqVhiDuv-VzkXx_W9trboMyzY"))
}

func TestVerifySignature(t *testing.T) {
	t.Parallel()
	receipt := []byte(`{"orderId":"GPA.xxxx-xxxx-xxxx-xxxxx","packageName":"my.package","productId":"myproduct","purchaseTime":1437564796303,"purchaseState":0,"developerPayload":"user001","purchaseToken":"some-token"}`)

	type in struct {
		pubkey  string
		receipt []byte
		sig     string
	}

	tests := []struct {
		name  string
		in    in
		err   error
		valid bool
	}{
		{
			name: "public key is invalid base64 format",
			in: in{
				pubkey:  "dummy_public_key",
				receipt: receipt,
				sig:     "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4=",
			},
			err:   errors.New("failed to decode public key"),
			valid: false,
		},
		{
			name: "public key is not rsa public key",
			in: in{
				pubkey:  "JTbngOdvBE0rfdOs3GeuBnPB+YEP1w/peM4VJbnVz+hN9Td25vPjAznX9YKTGQN4iDohZ07wtl+zYygIcpSCc2ozNZUs9pV0s5itayQo22aT5myJrQmkp94ZSGI2npDP4+FE6ZiF+7khl3qoE0rVZq4G2mfk5LIIyTPTSA4UvyQ=",
				receipt: receipt,
				sig:     "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4=",
			},
			err:   errors.New("failed to parse public key"),
			valid: false,
		},
		{
			name: "signature is invalid base64 format",
			in: in{
				pubkey:  "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB",
				receipt: receipt,
				sig:     "invalid_signature",
			},
			err:   errors.New("failed to decode signature"),
			valid: false,
		},
		{
			name: "signature is invalid",
			in: in{
				pubkey:  "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB",
				receipt: receipt,
				sig:     "JTbngOdvBE0rfdOs3GeuBnPB+YEP1w/peM4VJbnVz+hN9Td25vPjAznX9YKTGQN4iDohZ07wtl+zYygIcpSCc2ozNZUs9pV0s5itayQo22aT5myJrQmkp94ZSGI2npDP4+FE6ZiF+7khl3qoE0rVZq4G2mfk5LIIyTPTSA4UvyQ=",
			},
			err:   nil,
			valid: false,
		},
		{
			name: "normal",
			in: in{
				pubkey:  "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB",
				receipt: receipt,
				sig:     "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4=",
			},
			err:   nil,
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifySignature(tt.in.pubkey, tt.in.receipt, tt.in.sig)
			if err != nil {
				t.Error(err)
			}
			// if !reflect.DeepEqual(err, tt.err) {
			// 	t.Errorf("input: %v\nget: %s\nwant: %s\n", tt.in, err, tt.err)
			// }
		})
	}
}

// func TestRelaxedVictorAward(t *testing.T) {
// 	var awards int32
// 	var winerGroupNumber int
// 	var winerTotalNumber float64
// 	for i := 1; i <= 3; i++ {
// 		if i > int(math.Ceil(float64(3)*0.6)) {
// 			awards += 100
// 			continue
// 		}
// 		winerGroupNumber += i
// 	}
// 	for i := 1; i <= int(math.Ceil(float64(3)*0.6)); i++ {
// 		winerTotalNumber += float64(winerGroupNumber) / float64(i)
// 	}
// 	fmt.Println(awards, winerGroupNumber, winerTotalNumber)
// 	for i := 1; i <= 3 && i <= int(math.Ceil(float64(3)*0.6)); i++ {
// 		//x.GroupList[i-1].Awards = int32(math.Floor(float64(winerGroupNumber) / float64(i) / float64(winerGroupNumber) * float64(awards)))
// 		fmt.Println(int32(math.Floor(float64(winerGroupNumber) / float64(i) / float64(winerTotalNumber) * float64(awards))))
// 	}
// }
