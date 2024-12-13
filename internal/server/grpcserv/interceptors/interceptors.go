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
	pb "github.com/Ord1nI/metrix/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor{
	return grpc.UnaryServerInterceptor(func (ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any,error) {
		v, ok := req.(*pb.Metric)
		if !ok {
			return nil, errors.New("unknown type")
		}

		res, err := handler(ctx, req)

		switch q := res.(type) {
		case *pb.Error:
			log.Infoln(
				"\nREQ\n",
				v.String(),
				"\nRESPONSE\n",
				q.String(),
				"With error: ", err, "\n",
			)
		case *pb.Metric:
			log.Infoln(
				"\nREQ\n",
				v.String(),
				"\nRESPONSE\n",
				q.String(),
				"With error: ", err, "\n",
			)
		default:
			log.Info(
				"With error: ", err, "\n",
			)
		}

		return res, err

	})
}

func SignInterceptor(l logger.Logger, key []byte) grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func (ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any,error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("Cant get metadata")
		}

		stringHash := md["HashSHA256"][0]

		if stringHash != "" {
			if len(stringHash) < 64 {
				l.Infoln("Bad hash")
				return nil, errors.New("bad hash")
			}

			getHash, err := hex.DecodeString(stringHash)
			if err != nil {
				l.Infoln("error whiele decoding hex", err)
				return nil, errors.New("error whiele decoding hex")
			}


			bodyBytes := bytes.NewBuffer(nil)

			err = binary.Write(bodyBytes, binary.LittleEndian, req)

			signer := hmac.New(sha256.New, key)
			_, err = signer.Write(bodyBytes.Bytes())

			if err != nil {
				l.Infoln("Error while signing")
				return nil, errors.New("error while signing")
			}

			Hash := signer.Sum(nil)

			if !hmac.Equal(getHash, Hash) {
				l.Infoln("Hashes not equal")
				l.Infoln(getHash, "\n", Hash)
				return nil, errors.New("Hahs not equal")
			}

			return handler(ctx, req)
		} else {
			return handler(ctx, req)
		}
	})
}
