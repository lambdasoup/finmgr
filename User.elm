port module User exposing (..)

import Html exposing (..)
import Service exposing (User, Empty)


port getUserEmpty : Empty -> Cmd msg


port setUser : (User -> msg) -> Sub msg


getUser : Cmd msg
getUser =
    getUserEmpty <| Empty ""


type alias Model =
    Maybe User


init : ( Model, Cmd Msg )
init =
    ( Nothing, getUser )


type Msg
    = GetUser
    | SetUser User


view : Model -> Html Msg
view model =
    case model of
        Nothing ->
            text "loading user..."

        Just user ->
            text <| "user: " ++ user.email


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetUser ->
            ( Nothing, getUser )

        SetUser user ->
            ( Just user, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ setUser SetUser
        ]
