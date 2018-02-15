//# sourceURL=app.js
'use strict';

import elm from "../elm-src/Main.elm"

import * as grpc from "./grpc";
import * as push from "./push";

var app = elm.Main.fullscreen()
grpc.connect(app);
push.connect(app);
