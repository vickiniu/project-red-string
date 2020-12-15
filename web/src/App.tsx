import React from "react";
import "./App.css";
import { BrowserRouter, Route, Switch } from "react-router-dom";

import Home from "./components/Home";
import Individual from "./components/Individual";

function App() {
    return (
        <BrowserRouter>
            <Switch>
                <Route path="/individual/:individualID">
                    <Individual />
                </Route>
                <Route path="/">
                    <Home />
                </Route>
            </Switch>
        </BrowserRouter>
    );
}

export default App;
