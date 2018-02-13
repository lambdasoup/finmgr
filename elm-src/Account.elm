port module Account exposing (..)

import Html exposing (..)
import Service exposing (AccountInfo, Accounts)


port getAccounts : AccountInfo -> Cmd msg


port reply : (Accounts -> msg) -> Sub msg


type alias Model =
    { accountInfo : AccountInfo
    , accounts : Accounts
    }


init : ( Model, Cmd Msg )
init =
    ( { accountInfo = AccountInfo "" "" ""
      , accounts = Accounts ""
      }
    , Cmd.none
    )


type Msg
    = GetAccounts
    | ReplyReceived Accounts
    | SetAccountId String
    | SetPin String
    | SetBlz String


view : Model -> Html Msg
view model =
    text "loading user..."


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetAccounts ->
            ( model, getAccounts model.accountInfo )

        SetAccountId updated ->
            let
                updated1 =
                    model.accountInfo

                updated2 =
                    { updated1 | id = updated }
            in
                ( { model | accountInfo = updated2 }, Cmd.none )

        SetBlz updated ->
            let
                updated1 =
                    model.accountInfo

                updated2 =
                    { updated1 | blz = updated }
            in
                ( { model | accountInfo = updated2 }, Cmd.none )

        SetPin updated ->
            let
                updated1 =
                    model.accountInfo

                updated2 =
                    { updated1 | pin = updated }
            in
                ( { model | accountInfo = updated2 }, Cmd.none )

        ReplyReceived reply ->
            ( { model | accounts = reply }, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ reply ReplyReceived
        ]
