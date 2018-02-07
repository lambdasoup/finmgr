module Example exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: example.proto

import Protobuf exposing (..)

import Json.Decode as JD
import Json.Encode as JE


type alias HelloRequest =
    { name : String -- 1
    }


helloRequestDecoder : JD.Decoder HelloRequest
helloRequestDecoder =
    JD.lazy <| \_ -> decode HelloRequest
        |> required "name" JD.string ""


helloRequestEncoder : HelloRequest -> JE.Value
helloRequestEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "name" JE.string "" v.name)
        ]


type alias HelloReply =
    { message : String -- 1
    }


helloReplyDecoder : JD.Decoder HelloReply
helloReplyDecoder =
    JD.lazy <| \_ -> decode HelloReply
        |> required "message" JD.string ""


helloReplyEncoder : HelloReply -> JE.Value
helloReplyEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "message" JE.string "" v.message)
        ]
