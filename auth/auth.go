package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	authProto "github.com/cyruslo/proto/auth"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
)

//json pay load like {"exp":1617191803,"orig_iat":1617190303,"userid":"71875413"}
type authJWTPayLoad struct {
	PlayerID string  `json:"userid"`
}

func VerificationToken(token string) (verifiedID int64, err error) {
	md := metadata.Pairs("authorization", "Bearer "+token)
	// 新建一个有 metadata 的 context
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	em := &empty.Empty{}
	var rsp *authProto.VerifyResult
	authClient := authProto.DefaultClient()
	if authClient == nil {
		err = errors.New("create auth client fail")
		return
	}
	rsp, err = authClient.VerifyToken(ctx, em)
	if err != nil {
		err = errors.New(fmt.Sprintf("token:%s verification err:%s", token, err.Error()))
		return
	}
	if rsp.Verify == false {
		err = errors.New(fmt.Sprintf("token:%s verification fail", token))
		return
	}
	verifiedID, err = decodeJwt(token)
	return
}

func decodeJwt(token string) (pid int64, err error) {
	s := strings.Split(token, ".")
	payloadData := s[1]
	var decodedData []byte
	if decodedData, err = base64.RawURLEncoding.DecodeString(payloadData);err != nil {
		return
	}
	payLoad := &authJWTPayLoad{}
	if err = json.Unmarshal(decodedData, payLoad);err != nil {
		return
	}
	pid = stringToInt64(payLoad.PlayerID)
	if pid <= 0 {
		err = errors.New(fmt.Sprintf("decode jwt payload, pid invalid, payload string:%s", string(decodedData)))
		return
	}
	return
}

func stringToInt64(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return v
}