port module Main exposing (main)

import Html exposing (..)
import Service exposing (AccountInfo, Accounts)
import User exposing (..)


port getAccounts : AccountInfo -> Cmd msg


port reply : (Accounts -> msg) -> Sub msg


type alias Model =
    { accountInfo : AccountInfo
    , accounts : Accounts
    , userModel : User.Model
    }


init : ( Model, Cmd Msg )
init =
    let
        ( userInitModel, userInitCmd ) =
            User.init
    in
        ( { accountInfo = AccountInfo "" "" ""
          , accounts = Accounts ""
          , userModel = userInitModel
          }
        , Cmd.map UserMsg userInitCmd
        )


type Msg
    = UserMsg User.Msg
    | GetAccounts
    | ReplyReceived Accounts
    | SetAccountId String
    | SetPin String
    | SetBlz String


view : Model -> Html Msg
view model =
    Html.div []
        [ Html.map UserMsg (User.view model.userModel)
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetAccounts ->
            ( model, getAccounts model.accountInfo )

        SetAccountId updated ->
            let
                updated1 =
                    model.accountInfo

                updated2 =
                    { updated1 | id = updated }
            in
                ( { model | accountInfo = updated2 }, Cmd.none )

        SetBlz updated ->
            let
                updated1 =
                    model.accountInfo

                updated2 =
                    { updated1 | blz = updated }
            in
                ( { model | accountInfo = updated2 }, Cmd.none )

        SetPin updated ->
            let
                updated1 =
                    model.accountInfo

                updated2 =
                    { updated1 | pin = updated }
            in
                ( { model | accountInfo = updated2 }, Cmd.none )

        ReplyReceived reply ->
            ( { model | accounts = reply }, Cmd.none )

        UserMsg userMsg ->
            let
                ( updatedUserModel, userCmd ) =
                    User.update userMsg model.userModel
            in
                ( { model | userModel = updatedUserModel }, Cmd.map UserMsg userCmd )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ reply ReplyReceived
        , Sub.map UserMsg (User.subscriptions model.userModel)
        ]


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
