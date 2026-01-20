import i18n from "@/core/i18n";

export default i18n([
    [
        '密码策略',
        '密碼原則', 
        'Password Policy', 
    ],

    [
        '启用密码错误锁定：',
        '啟用密碼錯誤鎖定：', 
        'Lockdown Criteria', 
    ],

    [
        '用户在任意情况下，密码连续输错',
        '使用者在任何情況下，密碼連續輸錯', 
        'An account will be locked after'
    ],

    [
        '次，则账号被锁定',
        '次，則帳號會被鎖定', 
        'unsuccessful password attempts.', 
    ],

    [
        '账号被锁定后，',
        '帳戶被鎖定後，', 
        'The lock of an account will be undone in', 
    ],

    [
        '分钟后自动解锁或由管理员解锁',
        '分鐘後將自動解鎖或由管理員解鎖', 
        'minutes', 
    ],

    [
        '密码长度为10-32个字符，需同时包含大小写字母及数字 ',
        '密碼長度為10-32個字元，需同時包含大小寫字母及數位 ',
        'Password should contain numbers, letters in both uppercase and lowercase within 10~32 characters.', 
    ],

    [
        '密码有效期：',
        '密碼有效期間：', 
        'Expiration', 
    ],

    [
        '天',
        '天', 
        'day(s)',
    ],

    [
        '个月',
        '個月', 
        'month(s)', 
    ],

    [
        '永久有效',
        '永久有效', 
        'Permanent', 
    ],

    [
        '密码仅在指定时间段内有效，若超过该有效期，则需要修改密码，否则无法登录',
        '密碼僅在指定時間段內有效，若超過該有效期，則需要修改密碼，否則無法登入', 
        'Password is only valid within the specified period. You will need to change the password after it expires.', 
    ],

    [
        '保存',
        '儲存', 
        'Save', 
    ],

    [
        '保存成功',
        '儲存成功', 
        'Saved', 
    ],

    [
        '取消',
        '取消', 
        'Cancel', 
    ],

    [
        '密码强度：',
        '密碼強度：', 
        'Strength', 
    ],

    [
        '强密码',
        '強密碼', 
        'Strong', 
    ],

    [
        '弱密码',
        '弱密碼', 
        'Weak', 
    ],

    [
        '强密码格式：密码长度至少为',
        '強密碼格式: 密碼長度至少為', 
        'Principle: at least', 
    ],

    [
        '个字符，需同时包含 大小写英文字母、数字与特殊字符',
        '個字元，需同時包含 大小寫英文字母、數字與特殊字元', 
        'characters, and must contain numbers, letters in both uppercase and lowercase, and special characters', 
    ],

    [
        '弱密码格式：密码长度至少为6个字符',
        '弱密碼格式：密碼長度至少為6個字元', 
        'Principle: at least 6 characters', 
    ],

    [
        '强密码长度至少为${min}个字符。',
        '強密碼長度至少為${min}個字元', 
        'Strong Password should be at least ${min} characters.', 
    ],

    [
        '强密码长度最多为99个字符。',
        '強密碼長度最多為99個字元。', 
        'Strong Password should be at most 99 characters.', 
    ],

    [
        '密码错误次数范围为1~${max}。',
        '密碼錯誤次數範圍為1~${max}。', 
        'Password attempts should be within 1~${max}.', 
    ],

    [
        '锁定时间范围为10~180分钟。',
        '鎖定時間範圍為10~180分鐘。', 
        'Lockout time should be within 10~180 minutes.', 
    ],

    [
        '忘记密码重置：',
        '忘記密碼重置', 
        'Ways to Reset:', 
    ],

    [
        '通过短信验证',
        '通過短訊驗證', 
        'Via SMS', 
    ],

    [
        '通过邮箱验证',
        '通過郵箱驗證', 
        'Via Email', 
    ],

    [
        '用户忘记密码时，可以通过绑定的${phone}邮箱发送验证码验证身份，重新设置密码（管控密码的用户除外）',
        '使用者忘記密碼時，可以通過綁定的${phone}郵箱發送驗證碼驗證身份，重新設定密碼（管控密碼的使用者除外）', 
        'Passwords can be reset via the bound ${phone} Email except for unchangeable ones', 
    ],

    [
        '手机或',
        '手機或', 
        'phone or', 
    ],

    [
        '（短信验证需要配置对应的短信服务器插件才能生效，如果您还没有配置，可以在',
        '（短信驗證需要配置對應的短信伺服器插件才能生效，如果您還沒有配置，可以在', 
        "(SMS verification requires a SMS server plugin to take effect. If you haven't configured it, please complete it in", 
    ],

    [
        '（邮箱验证需要配置对应的邮箱服务器插件才能生效，如果您还没有配置，可以在',
        '（郵箱驗證需要配置對應的郵箱伺服器插件才能生效，如果您還沒有配置，可以在', 
        "(Email verification requires a Email server plugin to take effect. If you haven't configured it, please complete it in", 
    ],

    [
        ' 第三方消息插件 ',
        ' 第三方訊息插件 ', 
        ' Third-party Message Plugin ', 
    ],

    [
        ' 邮件服务 ',
        ' 邮件服务 ', 
        ' Email Service ', 
    ],

    [
        '页面进行操作）',
        '页面進行操作）', 
        'page.)', 
    ],

    [
        '设置 密码时效为 ${time}天 成功',
        '設定 密碼時效為 ${time}天 成功', 
        'Set password validity as ${time} day(s) succeed', 
    ],

    [
        '设置 密码时效为 永久有效 成功',
        '設定 密碼時效為 永久有效 成功', 
        'Set password validity as Permanent succeed', 
    ],

    [
        '设置 密码策略 成功',
        '設定 密碼原則 成功', 
        'Set password policy successfully', 
    ],

    [
        '启用 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
        '啟用 密碼錯誤鎖定，最大連續輸錯密碼次數為${passwdErrCnt}次，自動解鎖時間為${passwdLockTime}分鐘',
        'Enable password lock. Maximum password attempts: ${passwdErrCnt} times. Lock period: ${passwdLockTime} minutes.',
    ],

    [
        '关闭 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
        '關閉 密碼錯誤鎖定，最大連續輸錯密碼次數為${passwdErrCnt}次，自動解鎖時間為${passwdLockTime}分鐘', 
        'Disable password lock. Maximum password attempts: ${passwdErrCnt} times. Lock period: ${passwdLockTime} minutes.', 
    ],

    [
        '强密码，密码长度至少为${strongPwdLength}个字符',
        '強密碼，密碼長度至少為${strongPwdLength}個字元', 
        'Strong Password,Password should be at least ${strongPwdLength} characters.', 
    ],

    [
        '启用 忘记密码重置通过短信验证 成功',
        '啟用 忘記密碼重置通過短訊驗證 成功', 
        'Enable SMS Password Reset successfully', 
    ],

    [
        '关闭 忘记密码重置通过短信验证 成功',
        '關閉 忘記密碼重置通過短訊驗證 成功', 
        'Close SMS Password Reset successfully', 
    ],

    [
        '启用 忘记密码重置通过邮箱验证 成功',
        '啟用 忘記密碼重置通過郵箱驗證 成功', 
        'Enable Email Password Reset successfully', 
    ],

    [
        '关闭 忘记密码重置通过邮箱验证 成功',
        '關閉 忘記密碼重置通過郵箱驗證 成功', 
        'Close Email Password Reset successfully', 
    ],

    [
        '此项不允许为空。',
        '必填欄位。', 
        'This item cannot be empty.',
    ],

    [
        '请输入${min}-${max}的数值',
        '請輸入${min}-${max}的數值', 
        'Value from ${min} to ${max}', 
    ],

    [
        '用户初始密码设置：',
        '使用者初始密碼設定：', 
        'User Initial Password Setting:', 
    ],

    [
        '初始密码：',
        '初始密碼：', 
        'Initial Password:', 
    ],

    [
        '重置密码',
        '重設密碼', 
        'Reset', 
    ],

    [
        '随机密码',
        '隨機密碼', 
        'Random', 
    ],

    [
        '初始密码生成后无法再次查看，请妥善保管',
        '初始密碼生成後無法再次檢視，請妥善保管', 
        'Keep this password safe. After the initial password is generated, you will be unable to view it again', 
    ]],

);