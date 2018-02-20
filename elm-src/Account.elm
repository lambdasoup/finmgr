port module Account exposing (..)

import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Service exposing (AddBankRequest, Bank, BanksResponse, Account)


port getBanks : () -> Cmd msg


port addBank : AddBankRequest -> Cmd msg


port setBanks : (BanksResponse -> msg) -> Sub msg


type alias FormValues =
    { id : String
    , blz : String
    , pin : String
    }


type alias Model =
    { banks : List Bank
    , formValues : FormValues
    }


init : ( Model, Cmd Msg )
init =
    ( { banks = []
      , formValues = FormValues "" "" ""
      }
    , getBanks ()
    )


type Msg
    = GetAccounts
    | ReplyReceived BanksResponse
    | UpdateForm FormMsg


type FormMsg
    = SetAccountId String
    | SetPin String
    | SetBlz String
    | AddBank


view : Model -> Html Msg
view model =
    section []
        [ h1 [] [ text "Banks" ]
        , Html.map UpdateForm <| viewForm model.formValues
        , ul [] (List.map viewBank model.banks)
        ]


viewForm : FormValues -> Html FormMsg
viewForm values =
    div []
        [ input [ placeholder "account ID", onInput SetAccountId ] []
        , input [ placeholder "blz", onInput SetBlz ] []
        , input [ placeholder "pin", onInput SetPin ] []
        , button [ onClick AddBank ] [ text "add" ]
        ]


viewBank : Bank -> Html Msg
viewBank bank =
    li []
        [ h1 [] [ text bank.blz ]
        , text <|
            if bank.updating then
                "loading..."
            else
                "done."
        , ul [] (List.map viewAccount bank.accounts)
        ]


viewAccount : Account -> Html Msg
viewAccount account =
    li []
        [ text account.name
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetAccounts ->
            ( model, getBanks () )

        ReplyReceived reply ->
            ( { model | banks = reply.banks }, Cmd.none )

        UpdateForm formMsg ->
            let
                ( updatedFormValues, formCmd ) =
                    updateForm formMsg model.formValues
            in
                ( { model | formValues = updatedFormValues }, formCmd )


updateForm : FormMsg -> FormValues -> ( FormValues, Cmd Msg )
updateForm msg values =
    case msg of
        SetAccountId value ->
            ( { values | id = value }, Cmd.none )

        SetBlz value ->
            ( { values | blz = value }, Cmd.none )

        SetPin value ->
            ( { values | pin = value }, Cmd.none )

        AddBank ->
            ( values, addBank values )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ setBanks ReplyReceived
        ]
