port module Main exposing (main)

import Html exposing (Html, text, div, button)
import Html.Events exposing (onClick)
import User exposing (..)
import Account exposing (..)


type PushState
    = NotAvailable
    | Subscribed
    | Unsubscribed
    | NeedsPermission
    | Unknown
    | Invalid String


convert : String -> PushState
convert value =
    case value of
        "NotAvailable" ->
            NotAvailable

        "Unsubscribed" ->
            Unsubscribed

        "Subscribed" ->
            Subscribed

        "NeedsPermission" ->
            NeedsPermission

        _ ->
            Invalid value


port getPushState : () -> Cmd msg


port requestPermission : () -> Cmd msg


port subscribe : () -> Cmd msg


port unsubscribe : () -> Cmd msg


port setPushState : (String -> msg) -> Sub msg


type alias Model =
    { userModel : User.Model
    , accountModel : Account.Model
    , pushState : PushState
    }


init : ( Model, Cmd Msg )
init =
    let
        ( userInitModel, userInitCmd ) =
            User.init

        ( accountInitModel, accountInitCmd ) =
            Account.init
    in
        ( { userModel = userInitModel
          , accountModel = accountInitModel
          , pushState = Unknown
          }
        , Cmd.batch
            [ Cmd.map UserMsg userInitCmd
            , Cmd.map AccountMsg accountInitCmd
            , getPushState ()
            ]
        )


type Msg
    = UserMsg User.Msg
    | AccountMsg Account.Msg
    | SetPushState PushState
    | RequestPermission
    | Subscribe
    | Unsubscribe


view : Model -> Html Msg
view model =
    Html.div []
        [ viewPushState model.pushState
        , Html.map UserMsg <| User.view model.userModel
        , Html.map AccountMsg <| Account.view model.accountModel
        ]


viewPushState : PushState -> Html Msg
viewPushState state =
    Html.div []
        [ text <| toString state
        , case state of
            NeedsPermission ->
                button [ onClick RequestPermission ] [ text "request permission" ]

            Unsubscribed ->
                button [ onClick Subscribe ] [ text "subscribe" ]

            Subscribed ->
                button [ onClick Unsubscribe ] [ text "unsubscribe" ]

            _ ->
                text ""
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UserMsg userMsg ->
            let
                ( updatedUserModel, userCmd ) =
                    User.update userMsg model.userModel
            in
                ( { model | userModel = updatedUserModel }, Cmd.map UserMsg userCmd )

        AccountMsg accountMsg ->
            let
                ( updatedAccountModel, accountCmd ) =
                    Account.update accountMsg model.accountModel
            in
                ( { model | accountModel = updatedAccountModel }, Cmd.map AccountMsg accountCmd )

        SetPushState state ->
            ( { model | pushState = state }, Cmd.none )

        RequestPermission ->
            ( model, requestPermission () )

        Subscribe ->
            ( model, subscribe () )

        Unsubscribe ->
            ( model, unsubscribe () )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ Sub.map UserMsg (User.subscriptions model.userModel)
        , Sub.map AccountMsg (Account.subscriptions model.accountModel)
        , setPushState (\s -> SetPushState (convert s))
        ]


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
