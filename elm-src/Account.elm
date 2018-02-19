port module Account exposing (..)

import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Service exposing (Bank, Accounts)


port getAccounts : () -> Cmd msg


port addBank : Bank -> Cmd msg


port reply : (Accounts -> msg) -> Sub msg


type alias Model =
    { bank : Bank
    , accounts : Accounts
    }


init : ( Model, Cmd Msg )
init =
    ( { bank = Bank "" "" ""
      , accounts = Accounts [] False
      }
    , Cmd.none
    )


type Msg
    = GetAccounts
    | AddBank
    | ReplyReceived Accounts
    | SetAccountId String
    | SetPin String
    | SetBlz String


view : Model -> Html Msg
view model =
    section []
        [ h1 [] [ text "Accounts" ]
        , input [ placeholder "account ID", onInput SetAccountId ] []
        , input [ placeholder "blz", onInput SetBlz ] []
        , input [ placeholder "pin", onInput SetPin ] []
        , button [ onClick AddBank ] [ text "add" ]
        , text <|
            if model.accounts.loading then
                "loading..."
            else
                "done."
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetAccounts ->
            ( model, getAccounts () )

        AddBank ->
            ( model, addBank (model.bank) )

        SetAccountId updated ->
            let
                updated1 =
                    model.bank

                updated2 =
                    { updated1 | id = updated }
            in
                ( { model | bank = updated2 }, Cmd.none )

        SetBlz updated ->
            let
                updated1 =
                    model.bank

                updated2 =
                    { updated1 | blz = updated }
            in
                ( { model | bank = updated2 }, Cmd.none )

        SetPin updated ->
            let
                updated1 =
                    model.bank

                updated2 =
                    { updated1 | pin = updated }
            in
                ( { model | bank = updated2 }, Cmd.none )

        ReplyReceived reply ->
            ( { model | accounts = reply }, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ reply ReplyReceived
        ]
