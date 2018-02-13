port module User exposing (..)

import Html exposing (..)
import Navigation exposing (..)
import Html.Events exposing (..)
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
    | Logout String


view : Model -> Html Msg
view model =
    case model of
        Nothing ->
            text "loading user..."

        Just user ->
            userView user


userView : User -> Html Msg
userView user =
    div []
        [ text <| "user: " ++ user.email
        , button [ onClick <| Logout user.logoutUrl ] [ text "logout" ]
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetUser ->
            ( Nothing, getUser )

        SetUser user ->
            ( Just user, Cmd.none )

        Logout url ->
            ( model, Navigation.load url )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ setUser SetUser
        ]
