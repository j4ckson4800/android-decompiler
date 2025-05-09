# Example

Example command, for some apk with proto files included

```bash
proto-gen -apk com.komorebi.minimal.calendar.apk -o ./example
```

## Output example

Here's an output for the above command for package com/moloco/sdk/UserIntent.proto

```proto
syntax = "proto3";
package com.moloco.sdk.UserIntent;
option go_package = "com/moloco/sdk/UserIntent";


message UserAdInteractionExt {
    oneof info_ext {
        UserIntent.UserAdInteractionExt.ImpressionInteraction imp_interaction = 100;
        UserIntent.UserAdInteractionExt.ClickInteraction click_interaction = 101;
        UserIntent.UserAdInteractionExt.AppForegroundingInteraction app_foregrounding_interaction = 102;
        UserIntent.UserAdInteractionExt.AppBackgroundingInteraction app_backgrounding_interaction = 103;
    }
    string mref = 1;
    int64 client_timestamp = 2;
    string advertising_id = 3;
    UserIntent.UserAdInteractionExt.Device device = 4;
    UserIntent.UserAdInteractionExt.App app = 5;
    UserIntent.UserAdInteractionExt.Network network = 6;
    UserIntent.UserAdInteractionExt.MolocoSDK sdk = 7;
    message ImpressionInteraction {
    }
    message AppBackgroundingInteraction {
    }
    message ClickInteraction {
        UserIntent.UserAdInteractionExt.Position click_pos = 1;
        UserIntent.UserAdInteractionExt.Size screen_size = 2;
        UserIntent.UserAdInteractionExt.Position view_pos = 3;
        UserIntent.UserAdInteractionExt.Size view_size = 4;
        repeated string buttons = 5;
    }
    message AppForegroundingInteraction {
        int64 bg_ts_ms = 1;
    }
    message Device {
        int32 os = 1;
        string os_ver = 2;
        string model = 3;
        float screen_scale = 4;
    }
    message Network {
        int32 connection_type = 1;
        string carrier = 2;
    }
    message Size {
        float w = 1;
        float h = 2;
    }
    message Position {
        float x = 1;
        float y = 2;
    }
    message Button {
        int32 type = 1;
        UserIntent.UserAdInteractionExt.Position pos = 2;
        UserIntent.UserAdInteractionExt.Size size = 3;
    }
    message App {
        string id = 1;
        string ver = 2;
    }
    message MolocoSDK {
        string core_ver = 1;
        string adapter_ver = 2;
    }
}











````