port module Main exposing (main)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Bootstrap.Grid as Grid
import Bootstrap.Form as Form
import Bootstrap.Form.Input as Input
import Bootstrap.Form.Select as Select
import Bootstrap.Form.Checkbox as Checkbox
import Bootstrap.Form.Radio as Radio
import Bootstrap.Form.Textarea as Textarea
import Bootstrap.Form.Fieldset as Fieldset
import Bootstrap.Button as Button

import Service exposing (Hello, Bye)

port hello : Hello -> Cmd msg


port reply : (Bye -> msg) -> Sub msg


type alias Model =
  { accountId: String
  , reply: String
  }


initialModel : Model
initialModel =
  { accountId = ""
  , reply = "no reply yet"
  }

init : ( Model, Cmd Msg )
init =
  ( initialModel, Cmd.none )


type Msg
    = GetAccounts
    | ReplyReceived Bye
    | SetAccountId String


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
                    [ Form.label [for "account-pin"] [ text "Account PIN"]
                    , Input.password [ Input.id "account-pin" ]
                    ]
                , Button.button [ Button.primary ] [ text "Show Accounts" ]
                ]
            ]
          , Grid.col [] [ text model.reply ]
          ]
      ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetAccounts ->
            ( model, hello ( Hello model.accountId ) )

        SetAccountId updated ->
            ( { model | accountId = updated }, Cmd.none )

        ReplyReceived bye ->
            ( { model | reply = bye.name }, Cmd.none )


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
