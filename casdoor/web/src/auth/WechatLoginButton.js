// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {createButton} from "react-social-login-buttons";
import {StaticBaseUrl} from "../Setting";

function Icon({ width = 24, height = 24, color }) {
    return <img src={`${StaticBaseUrl}/buttons/wechat.svg`} />;
}

const config = {
    text: "Sign in with Wechat",
    icon: Icon,
    iconFormat: name => `fa fa-${name}`,
    style: {background: "rgb(0,194,80)"},
    activeStyle: {background: "rgb(0,158,64)"},
};

const WechatLoginButton = createButton(config);

export default WechatLoginButton;
