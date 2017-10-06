-- webhook filter for github/gitlab
function filter_github(body, branch)
    ref = "\"refs/heads/" .. branch .. "\""
    return string.find(body, ref) ~= nil
end
