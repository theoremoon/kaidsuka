module Main exposing (main)

import Browser
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (href, src)
import Url exposing (Url)
import Url.Builder


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlChange = UrlChanged
        , onUrlRequest = UrlRequested
        }



-- MODEL


type alias Model =
    { pagename : String
    , url : Url.Url
    , key : Nav.Key
    }



--- INIT


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init _ url key =
    ( Model (Url.toString url) url key, Cmd.none )



--- UPDATE


type Msg
    = UrlChanged Url
    | UrlRequested Browser.UrlRequest


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UrlRequested urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        UrlChanged url ->
            ( { model | url = url, pagename = Url.toString url }, Cmd.none )



--- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



--- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "SPA EXAMPLE"
    , body =
        [ div []
            [ h1 [] [ text "HELLO IM THEOREMOON" ]
            , img [ src (Url.Builder.absolute [ "assets", "theoldmoon0602.png" ] []) ] []
            , p [] [ a [ href (Url.Builder.absolute [] []) ] [ text "Index" ] ]
            , p [] [ a [ href (Url.Builder.absolute [ "hello" ] []) ] [ text "Hello" ] ]
            , p [] [ text model.pagename ]
            ]
        ]
    }
