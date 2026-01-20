import i18n from "@/core/i18n";

export default i18n([
    [
        '启用访问者网段限制',
        '啟用訪客網段限制', 
        'Impose network segment restrictions on users',
    ],

    [
        '开启后，被绑定的访问者只能在指定的网段内登录客户端，未绑定的访问者无法登录客户端。',
        '啟用後，被綁定的訪客只能在指定的網段內登入用戶端，未綁定的訪客無法登入用戶端。', 
        'After enabling, the bound visitors shall only log into Client on the bound network segment; the unbound visitors are unable to log in.',
    ],

    [
        '开启后，被绑定的访问者只能在指定的网段内登录客户端，未绑定的访问者登录客户端不受限制。',
        '啟用後，被綁定的訪客只能在指定的網段內登入用戶端，未綁定的訪客登入用戶端不受限制。', 
        'After enabling, the bound visitors shall only log into Client on the bound network segment; the unbound visitors are unrestricted.', 
    ],

    [
        '启用 访问者网段限制 成功',
        '啟用 訪客網段限制 成功', 
        'Enable login restriction on visitor segment successfully',
    ],

    [
        '关闭 访问者网段限制 成功',
        '關閉 訪客網段限制 成功', 
        'Disable login restriction on visitor segment successfully',
    ],

    [
        '开启后，被绑定的访问者只能在指定的网段内登录客户端，未绑定的访问者${canLogin}',
        '啟用後，被綁定的訪客只能在指定的網段內登入用戶端，未綁定的訪客${canLogin}', 
        'After enabling, the bound visitors shall only log into Client on the bound network segment; the unbound visitors ${canLogin}',
    ],

    [
        '无法登录客户端',
        '無法登入用戶端', 
        'are unable to log in', 
    ],

    [
        '登录客户端不受限制',
        '登入用戶端不受限制。', 
        'are unrestricted.', 
        
    ],

    [
        '此输入项不允许为空。',
        '此輸入項不能為空。', 
        'This item cannot be empty.', 
    ],

    [
        '网段名称不能包含 \\ / : * ? " < > | 特殊字符，长度不能超过128个字符。',
        '網路區段名稱不能包含 \\ / : * ? " < > | 特殊字元，長度不能超過128個字元。', 
        'Segment name cannot include special characters \\ / : * ? " < > |, the length cannot exceed 128 characters.', 
    ],

    [
        'IP地址格式形如xxx.xxx.xxx.xxx，每段必须是 0~255 之间的整数。',
        'IP位址格式形如xxx.xxx.xxx.xxx，每段必須是 0~255 之間的整數。', 
        'IP address should be in xxx.xxx.xxx.xxx format, and each segment should be an integer within 0~255.', 
    ],

    [
        '子网掩码格式形如xxx.xxx.xxx.xxx，每段必须是 0~255 之间的整数。',
        '子網路遮罩格式形如xxx.xxx.xxx.xxx，每段必須是 0~255 之間的整數。', 
        'Subnet Mask should be in xxx.xxx.xxx.xxx format, and each segment should be an integer within 0~255.', 
    ],

    [
        '非法的网段掩码参数。',
        '非法的網段遮罩參數。', 
        'Illegal parameters of subnet mask.', 
    ],

    [
        '请求参数错误。',
        '要求參數錯誤。', 
        'Request parameter error.', 
    ],

    [
        '请求出现错误。',
        '要求出現錯誤。', 
        'Request error。', 
    ],

    [
        '请求过多。',
        '太多要求。', 
        'Too many requests.', 
    ],

    [
        '不允许此操作。',
        '此動作不允許。', 
        'This action is not allowed.', 
    ],

    [
        '网段“${network}”已存在，无法重复添加。',
        '網段“${network}”已存在，無法重複新增。', 
        'Add failed. The segment "${network}" already exists.', 
    ],

    [
        '该网段已不存在。',
        '該網段已不存在。', 
        'This network does not exist.', 
    ],

    [
        '该网段名称已存在，请重新输入。',
        '該網路名稱已存在，請重新輸入。', 
        'This name already exists. Please re-enter.', 
    ],

    [
        '该策略不存在',
        '該原則不存在', 
        'Policy does not exist', 
    ],

    [
        '保存',
        '儲存', 
        'Save', 
    ],

    [
        '取消',
        '取消', 
        'Cancel',
    ],

]);
