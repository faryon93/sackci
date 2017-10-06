#  Webhooks
To automatically start a build when the source code repository changes, the most convenient way is to set up a webhook.
You have to specify at least on `trigger_tokens` per project to use this feature.

The webhook is executed when sending a POST request to `http://localhost/api/v1/project/<id>/trigger?token=<token>`

## Filter
Most of the time, webhooks contain a JSON body describing the push event. So there is a need to filter all trigger events by the configured branch.
To achieve this the `trigger_filter` option of each project configuration can be populated with a lua script.

This example script filters for Github and Gitlab webhooks.
```
function trigger_filter(body, branch)
    ref = "\"/ref/head/" .. branch .. "\""
    return string.find(body, ref) ~= nil
end
```