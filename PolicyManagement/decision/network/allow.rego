package network

import data.network_info.is_enabled
import data.network_info.no_policy_accessor_enabled
import data.network_info.ipv4
import data.network_info.ipv6

default accessible = false
default check_user_exist = false

# 开关关闭，允许
accessible = true {
    is_enabled = false
}

# 超级管理员不受网段限制,允许
accessible = true {
    input.accessor_id == "234562BE-88FF-4440-9BFF-447F139871A2"  
}

# admin管理员不受网段限制,允许
accessible = true {
    input.accessor_id == "266c6a42-6131-4d62-8f39-853e7093701c"  
}

# security管理员不受网段限制,允许
accessible = true {
    input.accessor_id == "4bb41612-a040-11e6-887d-005056920bea"  
}

# audit管理员不受网段限制,允许
accessible = true {
    input.accessor_id == "94752844-BDD0-4B9E-8927-1CA8D427E699"  
}

# 用户绑定的ip符合要求,允许
accessible = true {
    input.ip_type == "ipv4"
    ipv4.users[input.accessor_id]
    some i
    input.ip >= ipv4.users[input.accessor_id].nets[i].start_ip
    input.ip <= ipv4.users[input.accessor_id].nets[i].end_ip
}

# 用户所在部门绑定的ip符合要求，允许
accessible = true {
    input.ip_type == "ipv4"
    ipv4.users[input.accessor_id]
    dps := ipv4.users[input.accessor_id].departments
    some i,j
    input.ip >= ipv4.departments[dps[i]][j].start_ip
    input.ip <= ipv4.departments[dps[i]][j].end_ip
}

# 用户绑定的ip符合要求,允许
accessible = true {
    input.ip_type == "ipv6"
    ipv6.users[input.accessor_id]
    some i
    input.ip >= ipv6.users[input.accessor_id].nets[i].start_ip
    input.ip <= ipv6.users[input.accessor_id].nets[i].end_ip
}

# 用户所在部门绑定的ip符合要求，允许
accessible = true {
    input.ip_type == "ipv6"
    ipv6.users[input.accessor_id]
    dps := ipv6.users[input.accessor_id].departments
    some i,j
    input.ip >= ipv6.departments[dps[i]][j].start_ip
    input.ip <= ipv6.departments[dps[i]][j].end_ip
}

check_user_exist = true {
    ipv4.users[input.accessor_id]
}

check_user_exist = true {
    ipv6.users[input.accessor_id]
}

# 开关关闭并且访问者没配置策略
accessible = true {
    check_user_exist = false
    no_policy_accessor_enabled = true
}
