port module Main exposing (main)

import Html exposing (..)
import User exposing (..)
import Account exposing (..)


type alias Model =
    { userModel : User.Model
    , accountModel : Account.Model
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
          }
        , Cmd.batch
            [ Cmd.map UserMsg userInitCmd
            , Cmd.map AccountMsg accountInitCmd
            ]
        )


type Msg
    = UserMsg User.Msg
    | AccountMsg Account.Msg


view : Model -> Html Msg
view model =
    Html.div []
        [ Html.map UserMsg (User.view model.userModel)
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


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ Sub.map UserMsg (User.subscriptions model.userModel)
        , Sub.map AccountMsg (Account.subscriptions model.accountModel)
        ]


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }