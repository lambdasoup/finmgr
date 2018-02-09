port module Main exposing (main)

import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)


port hello : String -> Cmd msg


port reply : (String -> msg) -> Sub msg


type alias Model =
    String


init : ( Model, Cmd Msg )
init =
    ( "---", Cmd.none )


type Msg
    = Send
    | ReplyReceived String


view : Model -> Html Msg
view model =
    div []
        [ button [ onClick Send ] [ text "say hello" ]
        , div [] [ text (toString model) ]
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Send ->
            ( model, hello "Alice" )

        ReplyReceived name ->
            ( name , Cmd.none )


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
