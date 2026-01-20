import i18n from '@/core/i18n';

export default i18n([
    [
        'SMTP服务器设置',
        'SMTP伺服器設定', 
        'SMTP Server Settintgs', 
    ],

    [
        '邮件服务器（SMTP）:',
        '郵件伺服器（SMTP）:', 
        'Mail Server (SMTP)', 
    ],

    [
        '安全连接:',
        '安全連線:', 
        'Secure Connection', 
    ],

    [
        '端口:',
        '連結埠:', 
        'Port:', 
    ],

    [
        '说明',
        '說明', 
        'Description', 
    ],

    [
        '开启Open Relay，需要邮件服务器已支持Open Relay方能操作成功，开启后，邮箱验证不需要输入密码；关闭Open Relay，邮箱验证则需要输入密码。',
        '開啟Open Relay，需要郵件伺服器已支援Open Relay才能使操作成功，開啟後，郵箱驗證不需要輸入密碼；關閉Open Relay，郵箱驗證則需要輸入密碼。', 'Please ensure that Open Relay is supported by your mail server before turning it on, or your operation will fail. Once Open Relay is enabled, you can directly log into mailbox without password. If you disable Open Relay, you will need to input the password for mailbox login.', 
    ],

    [
        'Open Relay:',
        'Open Relay:', 
        'Open Relay:', 
    ],

    [
        '邮箱地址:',
        '電子郵箱地址:', 
        'Mail Address:', 
    ],

    [
        '邮箱密码:',
        '電子郵箱密碼:', 
        'Password:', 
    ],

    [
        '测试',
        '測試', 
        'Test', 
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

    [
        '测试中...',
        '測試中...', 
        'Testing', 
    ],

    [
        '无',
        '無', 
        'None', 
    ],

    [
        '提示',
        '提示', 
        'Tips', 
    ],

    [
        '此输入项不允许为空。',
        '此輸入項不允許為空。', 
        'This item cannot be empty.', 
    ],

    [
        'SMTP服务器名只能包含 英文、数字 及 @-_. 字符，长度范围 3~100 个字符，请重新输入。',
        'SMTP伺服器名稱只能包含英文、數字 及 @-_. 字元，長度範圍 3~100 個字元，請重新輸入。', 
        'Letters, numbers and @-_. only, length from 3 to 100. Please re-enter.', 
    ],

    [
        '请输入1~65535范围内的整数。',
        '請輸入1~65535之間的整數。', 
        'Please enter integer within 1~65535.', 
    ],

    [
        '邮箱地址只能包含 英文、数字 及 @-_. 字符，格式形如 XXX@XXX.XXX，长度范围 5~100 个字符，请重新输入。',
        '郵箱位址只能包含 英文、數字 及 @-_. 字元，格式形如 XXX@XXX.XXX，長度範圍 5~100 個字元，請重新輸入。', 
        'The email address can only contain English characters, numbers and @ -_ characters in XXX @ XXX.XXX format, and it should be between 5 and 100 characters. Please re-enter.', 
    ],

    [
        '测试连接成功，指定的服务器可用，您可以进入邮箱查看测试邮件。',
        '測試連線成功，指定的伺服器可用，您可以進入電子信箱檢視郵件。', 
        'Test connection succeeded. The server is available. You may check the test message in your mailbox.', 
    ],

    [
        '设置 SMTP服务器 成功',
        '設定 SMTP伺服器 成功', 
        'SMTP server was set successfully', 
    ],

    [
        'SMTP服务器认证失败，用户名或密码错误。',
        'SMTP伺服器認證失敗，使用者名稱或密碼錯誤。', 
        'SMTP server authentication failed. Bad user name or password.', 
    ],

    [
        '邮件服务器地址“${server}”；安全连接“${safeMode}”；端口“${port}”；Open Relay“关闭”；邮箱地址“${email}”；邮箱密码“******”',
        '郵件伺服器位址“${server}”；安全連線“${safeMode}”；連接埠“${port}”；Open Relay“關閉”；電子郵件地址“${email}”；電子郵件密碼“******”', 
        'Mail server "${server}"; secure connection "${safeMode}"; port "${port}"; Open Relay "off"; mail address "${email}"; password "******"',
    ],

    [
        '邮件服务器地址“${server}”；安全连接“${safeMode}”；端口“${port}”；Open Relay“开启”；邮箱地址“${email}”',
        '郵件伺服器位址“${server}”；安全連線“${safeMode}”；連接埠“${port}”；Open Relay“開啟”；電子郵件地址“${email}”', 
        'Mail server "${server}"; secure connection "${safeMode}"; port "${port}"; Open Relay "on"; mail address "${email}"',
    ],

    [
        'SMTP服务器不可用，请检查服务器地址、安全连接或端口是否正确。',
        'SMTP伺服器不可用，請檢查伺服器位址、安全連線或連接埠是否正確。', 
        'SMTP server is unavailable, please check Mail server, secure connection or port', 
    ],

    [
        '测试连接失败，邮箱地址不正确，请重新输入。',
        '測試連線失敗，郵箱位址錯誤，請重新輸入。', 
        'Test connection failed. Invalid email address, please re-enter.', 
    ],

    [
        '测试连接失败，邮箱地址或密码不正确，请重新输入。',
        '測試連線失敗，郵箱位址或密碼錯誤，請重新輸入。', 
        'Test connection failed. Invalid email address or password, please re-enter.', 
    ],
]

);