port module Main exposing (main)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Bootstrap.Grid as Grid
import Bootstrap.Form as Form
import Bootstrap.Form.Input as Input
import Bootstrap.Button as Button

import Service exposing (AccountInfo, Accounts)

port getAccounts : AccountInfo -> Cmd msg


port reply : (Accounts -> msg) -> Sub msg


type alias Model =
  { accountInfo: AccountInfo
  , accounts: Accounts
  }


initialModel : Model
initialModel =
  { accountInfo = AccountInfo "" "" ""
  , accounts = Accounts ""
  }



init : ( Model, Cmd Msg )
init =
  ( initialModel, Cmd.none )


type Msg
    = GetAccounts
    | ReplyReceived Accounts
    | SetAccountId String
    | SetPin String
    | SetBlz String


view : Model -> Html Msg
view model =
  Grid.container []                                     -- Creates a div that centers content
      [ Grid.row []                                     -- Creates a row with no options
          [ Grid.col []
            [ Form.form [ onSubmit GetAccounts ]
                [ Form.group []
                    [ Form.label [for "account-id"] [ text "Account ID"]
                    , Input.text [ Input.onInput SetAccountId ]
                    ]
                , Form.group []
                    [ Form.label [for "account-blz"] [ text "Account BLZ"]
                    , Input.text [ Input.onInput SetBlz ]
                    ]
                , Form.group []
                    [ Form.label [for "account-pin"] [ text "Account PIN"]
                    , Input.password [ Input.onInput SetPin ]
                    ]
                , Button.button [ Button.primary ] [ text "Show Accounts" ]
                ]
            ]
          , Grid.col [] [ text model.accounts.info ]
          ]
      ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetAccounts ->
            ( model, getAccounts model.accountInfo )

        SetAccountId updated ->
            let
              updated1 = model.accountInfo
              updated2 = { updated1 | id = updated }
            in
              ( { model | accountInfo = updated2 }, Cmd.none )

        SetBlz updated ->
            let
              updated1 = model.accountInfo
              updated2 = { updated1 | blz = updated }
            in
              ( { model | accountInfo = updated2 }, Cmd.none )

        SetPin updated ->
            let
              updated1 = model.accountInfo
              updated2 = { updated1 | pin = updated }
            in
              ( { model | accountInfo = updated2 }, Cmd.none )

        ReplyReceived reply ->
            ( { model | accounts = reply }, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions model =
    reply ReplyReceived


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
