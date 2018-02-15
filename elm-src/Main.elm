port module Main exposing (main)

import Html exposing (..)
import User exposing (..)
import Account exposing (..)


type PushState
    = NotAvailable
    | Available
    | Unknown
    | Invalid String


convert : String -> PushState
convert value =
    case value of
        "NotAvailable" ->
            NotAvailable

        "Available" ->
            Available

        _ ->
            Invalid value


port getPushState : () -> Cmd msg


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
    | SetPushState String


view : Model -> Html Msg
view model =
    Html.div []
        [ text <| toString model.pushState
        , Html.map UserMsg <| User.view model.userModel
        , Html.map AccountMsg <| Account.view model.accountModel
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

        SetPushState str ->
            ( { model | pushState = convert str }, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ Sub.map UserMsg (User.subscriptions model.userModel)
        , Sub.map AccountMsg (Account.subscriptions model.accountModel)
        , setPushState SetPushState
        ]


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
