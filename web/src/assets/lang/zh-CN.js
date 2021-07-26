export const lang = {
  error: {
    text: "遇到一个致命错误，请将输出的信息发送给管理员"
  },
  form: {
    input: [
      {
        prepend: "语言",
      },
      {
        prepend: "密码",
        placeholder: "无需设置密码请留空"
      },
      {
        prepend: "过期时间"
      }
    ],
    textarea: {
      placeholder: {
        code: "Talk is cheap, Show me the code.",
        read_once: "阅后即焚"
      }
    },
    select: {
      plain: "纯文本",
      none: "无",
      hour: "一小时",
      day: "一天",
      week: "一周",
      month: "一个月",
      year: "一年"
    },
    submit: "保存",
    checkbox: "阅后即焚"
  },
  success: {
    h2: "保存成功",
    p: [
      {
        button: "返回主页"
      }
    ],
    ul: {
      li: [
        {
          text: "在导航栏中输入<strong>索引</strong>"
        },
        {
          browser: "在浏览器中访问",
          tooltip: "在新页面中查看"
        },
        {
          scan_qr_code: "扫描下方二维码"
        }
      ]
    },
    badge: {
      copy: "复制链接",
      success: "复制成功",
      fail: "复制失败"
    }
  },
  auth: {
    form: {
      label: "此 Paste 已加密，请输入密码：",
      button: "提交",
      placeholder: "密码错误"
    }
  },
  nav: {
    router_link: "返回主页",
    form: {
      placeholder: "索引",
      button: "前往"
    },
    beg: "给个 Star 吧 ~"
  },
  not_found: {
    content: {
      title: "您访问的页面没有找到",
      go_home: "返回主页"
    }
  },
  view: {
    parsed: "渲染",
    raw: "源码",
    lines: "行",
    lang: {
      cpp: "C/C++",
      java: "Java",
      bash: "Bash",
      html: "HTML",
      python: "Python",
      markdown: "Markdown",
      go: "Go",
      json: "JSON",
      plaintext: "纯文本"
    },
    copy: "复制",
    tooltip: {
      click: "点按以复制",
      success: "成功",
      fail: "失败"
    }
  }
};
