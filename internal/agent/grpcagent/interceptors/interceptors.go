package interceptors

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/Ord1nI/metrix/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SignInterceptor(l logger.Logger, key []byte) grpc.UnaryClientInterceptor {
	return grpc.UnaryClientInterceptor(func (ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		bodyBytes := bytes.NewBuffer(nil)

		err := binary.Write(bodyBytes, binary.LittleEndian, req)

		signer := hmac.New(sha256.New, key)
		_, err = signer.Write(bodyBytes.Bytes())

		if err != nil {
			l.Infoln("Error while signing")
			return errors.New("error while signing")
		}

		Hash := signer.Sum(nil)
		HashStr := hex.EncodeToString(Hash)

		ctx = metadata.AppendToOutgoingContext(ctx, "HashSHA256", HashStr)

		invoker(ctx, method, req, reply, cc)

		return nil
	})
}

func AddIPInterceptro(l logger.Logger, IP string) grpc.UnaryClientInterceptor {
	return grpc.UnaryClientInterceptor(func (ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		ctx = metadata.AppendToOutgoingContext(ctx, "X-Real-IP", IP)

		invoker(ctx, method, req, reply, cc)

		return nil
	})
}
