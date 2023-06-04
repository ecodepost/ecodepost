package slate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var htmlStr = `<p>“干得比驴累，吃得比猪差，起得比鸡早，睡得比狗晚，看上去比谁都好，五年后比谁都老。“——先自嘲下！
<img src="https://cdn.gocn.vip/ava/15.png" alt="" /></p>

<p>我们都知道，未来是互联网科技从业者的，更确切的说是属于程序员、设计师、产品经理的&hellip;</p>

<p>接私活请戳：<a href="http://www.mayigeek.com/">http://www.mayigeek.com/</a></p>`

/**
[{"type":"p","children":[{"text":"“干得比驴累，吃得比猪差，起得比鸡早，睡得比狗晚，看上去比谁都好，五年后比谁都老。“——先自嘲下！\n"},{"type":"img","children":[{}],"url":"http://ww1.sinaimg.cn/mw1024/8db48fb7jw1f8ro353i3kj209a09at9y.jpg"}]},{"text":"\n\n"},{"type":"p","children":[{"text":"我们都知道，未来是互联网科技从业者的，更确切的说是属于程序员、设计师、产品经理的…"}]},{"text":"\n\n"},{"type":"p","children":[{"text":"接私活请戳："},{"type":"a","children":[{"text":"http://www.mayigeek.com/"}]}]}]
*/
/**
[{"type":"p","children":[{"text":"“干得比驴累，吃得比猪差，起得比鸡早，睡得比狗晚，看上去比谁都好，五年后比谁都老。“——先自嘲下！\n"},{"type":"img","children":[{"text":""}],"url":"https://cdn.gocn.vip/ava/15.png"}]},{"text":"\n\n"},{"type":"p","children":[{"text":"我们都知道，未来是互联网科技从业者的，更确切的说是属于程序员、设计师、产品经理的…"}]},{"text":"\n\n"},{"type":"p","children":[{"text":"接私活请戳："},{"type":"a","url":"http://www.mayigeek.com/","children":[{"text":"http://www.mayigeek.com/"}]}]}]
*/
func TestHtmlToSlateJson(t *testing.T) {
	info, err := HtmlToSlateJson(htmlStr)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}

var htmlStrCode = `
<p>gops 是一个用来列出系统中正在使用的Go程序，同时还可以诊断正在运行的Go程序。</p>

<pre><code class="language-go">$ gops
983     uplink-soecks   (/usr/local/bin/uplink-soecks)
52697   gops    (/Users/jbd/bin/gops)
51130   gocode  (/Users/jbd/bin/gocode)
</code></pre>`

func TestHtmlToSlateJson2(t *testing.T) {
	info, err := HtmlToSlateJson(htmlStrCode)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}

var htmlStrCloseTag = `
<p>gops 是一个用来列出系统中正在使用的Go程序，同时还可以诊断正在运行的Go程序。</p>

<hr/>`

func TestHtmlToSlateJson3(t *testing.T) {
	info, err := HtmlToSlateJson(htmlStrCloseTag)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}

var newImageTag = `
<p>GoCN 本期老司机系列我们请来了游戏后端开发的达人，之前也分享过很多游戏开发经验 —— @达达 为大家解答关于Go 游戏开发方面的问题。

达达是来自真有趣信息科技有限公司的CTO，非科班老司机，在游戏和互联网行业摸爬滚打十余载，擅长各种降低研发难度和成本的奇技淫巧，曾负责《神仙道》、《仙侠道》等项目的服务端架构设计和研发，目前正在逐步整理并开源真有趣团队的游戏服务端技术架构和开发流程，请关注 <a href="https://github.com/funny/">https://github.com/funny/</a> 了解最新进展，欢迎广大同行参与开源框架的研发。

<img src="https://github.com/gocn/images/blob/master/269e6c88f90b5ceb69095566bd4503b5.jpeg?raw=true" alt=""/> <div />

不欢迎任何与主题无关的讨论和喷子。

下面欢迎大家对Go游戏开发方面的问题向 @达达 提问，请直接回帖提问！

&amp;gt; 达达的个人博客：<a href="http://1234n.com/">http://1234n.com/</a>
&amp;gt; 知乎的专栏：<a href="https://zhuanlan.zhihu.com/idada">https://zhuanlan.zhihu.com/idada</a>
&amp;gt; 最热门的一个回答：Go 的垃圾回收机制在实践中有哪些需要注意的地方？<a href="https://www.zhihu.com/question/21615032/answer/18781477">https://www.zhihu.com/question/21615032/answer/18781477</a>

本次活动持续两天，2016-10-17至2016-10-18，达达会在空闲时间上来给大家一一回答。</p>
`

func TestNewImageTagJson(t *testing.T) {
	info, err := HtmlToSlateJson(newImageTag)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}

var newPTag = `
<p><span>aa</span</p>
`

func TestPson(t *testing.T) {
	info, err := HtmlToSlateJson(newPTag)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}

var newPTag2 = `<h1>前言</h1>

<p>Android上抓包HTTPS是不是越来越难了？高版本无法添加CA证书，抓包软件依赖太多，VPN模式、或HOOK程序时，会被APP检测到。对抗成本愈加增高。有什么万能的工具吗？</p>

<p>是的，<a href="https://ecapture.cc">eCapture for Android</a>来了。以后在Android上抓HTTPS通讯包，再也不用安装CA证书了，再也不用下载一堆python依赖环境了，再也不用重打包ssl类库了，再也不用改一堆手机参数了，一键启用，简单明了。</p>

<h1>eCapture简介</h1>

<p>eCapture是一款无需CA证书即可抓获HTTPS明文的软件。支持pcapng格式，支持Wireshark直接查看。基于eBPF技术，仅需root权限，即可一键抓包。
eCapture中文名旁观者，即 当局者迷，旁观者清。
<img src="https://www.cnxct.com/wp-content/uploads/2022/09/ecapture.cc_zh_.png" alt="" /></p>

<p>2022年年初上海疫情期间，<a href="https://www.cnxct.com/what-is-ecapture/">笔者开始编写并开源</a>，至今已经半年，GitHub上已经 <code>4200</code>个星星。</p>

<p>eCapture是基于eBPF技术实现的抓包软件，依赖系统内核是否支持eBPF。目前支持在操作系统上，支持了X86_64\ARM64的Linux kernel 4.18以上内核，支持了ARM64 Android(Linux) kernel 5.4以上版本。最新版是在2022年9月9日发布的v0.4.3版本。</p>

<p><img src="https://www.cnxct.com/wp-content/uploads/2022/09/github.com_ehids_ecapture_releases.png" alt="" /></p>

<h1>演示视频</h1>

<p>下载后，一条命令启动，干净利索。
<code>./ecapture tls -w ecapture.pcapng</code></p>

<p>先看演示视频，演示环境为Ubuntu 21.04、Android 5.4 （Pixel 6）。</p>

<p><a href="https://www.bilibili.com/video/BV1xP4y1Z7HB/"><img src="https://cdn.gocn.vip/forum-user-images/20221008/03b7037e3f6241778d374272923519dc.jpg" alt="ecapture-tmp-play.jpg" /></a></p>

<h2>模块功能</h2>

<p>eCapture支持tls、bash、mysqld、postgres等模块的信息提取与捕获。本文仅讨论<code>tls</code>这个HTTPS/TLS明文捕获模块。</p>

<h2>加密通讯明文捕获&ndash;tls模块</h2>

<p><code>tls</code>模块在加密通讯类库上，支持了openssl、gnutls、nspr/nss、boringssl等类库。
但在Android上，pcapng模式只支持boringssl，文本模式则都支持。</p>

<h1>如何使用eCapture</h1>

<h2>环境依赖</h2>

<ol>
<li>操作系统 Linux kernel 4.18以上，Android kernel 5.4以上。</li>
<li>支持BPF，可选支持BTF（eCapture版本不同）</li>
<li>root权限</li>
</ol>

<h2>版本选择</h2>

<p>BPF <a href="https://facebookmicrosites.github.io/bpf/blog/2020/02/19/bpf-portability-and-co-re.html">CO-RE</a>特性为BTF通用的格式，用于做跨内核版本兼容。有的Android手机是没开启BTF。可以查看系统配置确认，<code>CONFIG_DEBUG_INFO_BTF=y</code>为开启；<code>CONFIG_DEBUG_INFO_BTF=n</code>为关闭；其他为不支持BPF，无法使用eCapture。</p>

<pre><code class="language-shell">cfc4n@vm-server:~$# cat /boot/config- | grep CONFIG_DEBUG_INFO_BTF
CONFIG_DEBUG_INFO_BTF=y
</code></pre>

<p>Android系统上，config是gzip压缩的，且配置文件目录也变了。可使用<code>zcat /proc/config.gz</code>命令代替。</p>

<p>eCapture默认发行了支持CO-RE的ELF程序。Android版会发行一个5.4内核不支持BTF（即没有CO-RE）的版本。 下载后，可以通过<code>./ecapture -v</code>确认。 非CO-RE版本的version信息中包含编译时的内核版本。</p>

<pre><code class="language-shel"># no CO-RE
eCapture version:	linux_aarch64:0.4.2-20220906-fb34467:5.4.0-104-generic
# CO-RE
eCapture version:	linux_aarch64:0.4.2-20220906-fb34467:[CORE]
</code></pre>

<p>若版本不符合自己需求，可以自行编译，步骤见文末。</p>

<h2>全局参数介绍</h2>

<p>全局参数重点看如下几个</p>

<pre><code class="language-shell">root@vm-server:/home/cfc4n/# ecapture -h
      --hex[=false]		print byte strings as hex encoded strings
  -l, --log-file=&quot;&quot;		-l save the packets to file
  -p, --pid=0			if pid is 0 then we target all pids
  -u, --uid=0			if uid is 0 then we target all users
</code></pre>

<ol>
<li><code>--hex</code> 用于stdout输出场景，展示结果的十六进制，用于查看非ASCII字符，在内容加密、编码的场景特别有必要。</li>
<li><code>-l, --log-file=</code>  保存结果的文件路径。</li>
<li><code>-p, --pid=0</code> 捕获的目标进程，默认为0，则捕获所有进程。</li>
<li><code>-u, --uid=0</code> 捕获的目标用户，默认为0，则捕获所有用户，对Android来说，是很需要的参数。</li>
</ol>

<h2>模块参数</h2>

<pre><code class="language-shell">root@vm-server:/home/cfc4n/project/ssldump# bin/ecapture tls -h
OPTIONS:
      --curl=&quot;&quot;		curl or wget file path, use to dectet openssl.so path, default:/usr/bin/curl
      --firefox=&quot;&quot;	firefox file path, default: /usr/lib/firefox/firefox.
      --gnutls=&quot;&quot;	libgnutls.so file path, will automatically find it from curl default.
      --gobin=&quot;&quot;	path to binary built with Go toolchain.
  -h, --help[=false]	help for tls
  -i, --ifname=&quot;&quot;	(TC Classifier) Interface name on which the probe will be attached.
      --libssl=&quot;&quot;	libssl.so file path, will automatically find it from curl default.
      --nspr=&quot;&quot;		libnspr44.so file path, will automatically find it from curl default.
      --port=443	port number to capture, default:443.
      --pthread=&quot;&quot;	libpthread.so file path, use to hook connect to capture socket FD.will automatically find it from curl.
      --wget=&quot;&quot;		wget file path, default: /usr/bin/wget.
  -w, --write=&quot;&quot;	write the  raw packets to file as pcapng format.
</code></pre>

<h4>-i参数</h4>

<p><code>-i</code>参数为网卡的名字，Linux上默认为<code>eth0</code>，Android上默认为<code>wlan0</code>，你可以用这个参数自行指定。</p>

<h2>输出模式</h2>

<p>输出格式支持两种格式，文本跟pcapng文件。有三个参数，</p>

<ol>
<li>默认，全局参数，输出文本结果到stdout</li>
<li><code>-l</code> 全局参数，保存文本结果的文件路径</li>
<li><code>-w</code> 仅tls模块参数，保存pcapng结果的文件路径</li>
</ol>

<h2>类库路径</h2>

<p>Linux上支持多种类路，不同类路的路径也不一样。</p>

<table>
<thead>
<tr>
<th>类库</th>
<th>参数路径</th>
<th>默认值</th>
</tr>
</thead>

<tbody>
<tr>
<td>openssl/boringssl</td>
<td>&ndash;libssl</td>
<td>Linux自动查找，Android为/apex/com.android.conscrypt/lib64/libssl.so</td>
</tr>

<tr>
<td>gnutls</td>
<td>&ndash;gnutls</td>
<td>Linux自动查找，Android pcapng模式暂未支持</td>
</tr>

<tr>
<td>nspr/nss</td>
<td>&ndash;nspr</td>
<td>Linux自动查找，Android pcapng模式暂未支持</td>
</tr>
</tbody>
</table>

<h3>文本模式</h3>

<p><code>-l</code> 或者不加 <code>-w</code> 参数将启用该模式。</p>

<p>支持openssl、boringssl、gnutls、nspr/nss等多种TLS加密类库。
支持DTLS、TLS1.0至TLS1.3等所有版本的加密协议。
支持<code>-p</code>、<code>-u</code>等所有全局过滤参数。</p>

<h3>pcapng模式</h3>

<p><code>-w</code> 参数启用该模式，并用<code>-i</code>选择网卡名，Linux系统默认为<code>eth0</code>,Android系统默认为<code>wlan0</code>，</p>

<p>仅支持openssl、boringssl两个类库的数据捕获。暂不支持TLS 1.3协议。</p>

<h3>类库与参数支持</h3>

<p>在Linux系统上，大部分类库与参数都是可以支持的。但在Android系统上，因为内核与ARM架构的原因，支持的参数上，有一定的差异。</p>

<h4>不同模式的参数支持</h4>

<p><code>-p</code>、<code>-u</code>两个全局参数，支持文本模式，不支持pcapng模式。这是因为pcapng模式是使用eBPF TC技术实现。
<img src="https://www.cnxct.com/wp-content/uploads/2022/09/how-ecapture-works.png" alt="" /></p>

<p>|  模式  | -p | -u |
|  &mdash;  | &mdash; | &mdash; |
|  文本  | ✅ | ✅ |
| pcapng | ❌ |❌  |</p>

<h4>不同模式的类库以协议支持</h4>

<table>
<thead>
<tr>
<th>模式</th>
<th>openssl（类库）</th>
<th>boringssl（类库）</th>
<th>TLS 1.0/1.<sup>1</sup>&frasl;<sub>1</sub>.2（协议）</th>
<th>TLS 1.3（协议）</th>
</tr>
</thead>

<tbody>
<tr>
<td>文本</td>
<td>✅</td>
<td>✅</td>
<td>✅</td>
<td>✅</td>
</tr>

<tr>
<td>pcapng</td>
<td>✅</td>
<td>✅</td>
<td>✅</td>
<td>❌</td>
</tr>
</tbody>
</table>
<p>pcapng模式暂时不支持TLS 1.3，<a href="https://github.com/ehids/ecapture/pull/143">TLS 1.3密钥捕获功能</a>已经开发完成，只是遇到一些BUG，还在解决中。 笔者不是openssl的专家，对TLS 协议也不太熟。需要补充这两块的知识，解决起来成本比较高，也欢迎对这块擅长的朋友一起来解决。</p>

<h2>与tcpdump联合使用</h2>

<p>eCapture基于eBPF TC，实现了流量捕获，并保存到pcapng文件中。基于eBPF Uprobe实现了TLS Master Secret的捕获。并基于Wireshark的<a href="https://github.com/pcapng/pcapng/pull/54">Decryption Secrets Block (DSB)</a>标准，实现了<a href="https://github.com/google/gopacket/pull/1042">gopacket的DSB功能</a>，合并网络包与密钥，保存到pcapng中。</p>

<p>eCapture在网络包捕获上，没有tcpdump强大，不支持丰富的参数。你可以用eCapture捕获master secrets，用tcpdump捕获网络包，然后使用wiresahrk自定义设置密钥文件，配合使用。</p>

<p><img src="https://image.cnxct.com/2022/09/ecapture-wireshark-scaled.jpg" alt="" /></p>

<h4>网络包捕获</h4>

<p>tcpdump 的常规用法，不再赘述。</p>

<h4>密钥捕获</h4>

<p>同时启用ecapture ，模式可以选文本或者pcapng，都会但会保存TLS的master secrets密钥数据到ecapture_masterkey.log中。</p>

<h4>网络包查看</h4>

<p>用Wireshark打开网络包文件，设置这个master key文件，之后就可以看到TLS解密后的明文了。</p>

<p>配置路径：<code>Wireshark</code> &ndash;&gt; <code>Preferences</code> &ndash;&gt; <code>Protocols</code> &ndash;&gt; <code>TLS</code> &ndash;&gt; <code>(Pre)-Master-Secret log filename</code>
<img src="https://www.cnxct.com/wp-content/uploads/2022/09/wireshark-master-secrets.jpg" alt="" /></p>

<h2>参数</h2>

<h3>指定路径</h3>

<h4>默认路径</h4>

<p>在Android上，Google使用了boring ssl类库，也就是C++语言在<code>libssl</code>基础上的包装。默认情况下，会使用<code>/apex/com.android.conscrypt/lib64/libssl.so</code>路径。</p>

<h4>APP的类库确认</h4>

<p>你可以使用<code>lsof -p {APP PID}|grep libssl</code>来确认。 若不是默认路径，则可以使用<code>--libssl</code>参数来指定。</p>

<h1>高级用法</h1>

<p>如果你需要查看的APP是自定义SSL类库，那么你可以自助修改eCapture来实现。</p>

<h3>自定义函数名与offset</h3>

<p>首先，需要确定HOOK函数的函数名或者符号表地址。</p>

<h4>没有源码</h4>

<p>如果你没有类库源码可以通过IDA等软件静态分析、动态调试，确定SSL Write的地址offset。在配置填写在<code>user/module/probe_openssl.go</code>文件中，对应的probe配置部分。</p>

<pre><code class="language-go">{
    Section:          &quot;uprobe/SSL_write&quot;,
    EbpfFuncName:     &quot;probe_entry_SSL_write&quot;,
    AttachToFuncName: &quot;SSL_write&quot;,
    UprobeOffset:       0xFFFF00, // TODO
    BinaryPath:       binaryPath,
},
</code></pre>

<h3>offset自动计算</h3>

<p>如果你有源码，则可以通过<code>offsetof</code>宏来自动计算。</p>

<pre><code class="language-c++">//  g++ -I include/ -I src/ ./src/offset.c -o off
#include &lt;stdio.h&gt;
#include &lt;stddef.h&gt;
#include &lt;ssl/internal.h&gt;
#include &lt;openssl/base.h&gt;
#include &lt;openssl/crypto.h&gt;

#define SSL_STRUCT_OFFSETS               \
    X(ssl_st, session)              \
    X(ssl_st, s3)              \
    X(ssl_session_st, secret)        \
    X(ssl_session_st, secret_length)  \
    X(bssl::SSL3_STATE, client_random) \
    X(bssl::SSL_HANDSHAKE, new_session) \
    X(bssl::SSL_HANDSHAKE, early_session) \
    X(bssl::SSL3_STATE, hs) \
    X(bssl::SSL3_STATE, established_session) \
    X(bssl::SSL_HANDSHAKE, expected_client_finished_)


struct offset_test
{
    /* data */
    int t1;
    bssl::UniquePtr&lt;SSL_SESSION&gt; session;
};

int main() {
    printf(&quot;typedef struct ssl_offsets { // DEF \n&quot;);
#define X(struct_name, field_name) \
    printf(&quot;   int &quot; #struct_name &quot;_&quot; #field_name &quot;; // DEF\n&quot;);
    SSL_STRUCT_OFFSETS
#undef X
    printf(&quot;} ssl_offsets; // DEF\n\n&quot;);

    printf(&quot;/* %s */\nssl_offsets openssl_offset_%d = { \n&quot;,
           OPENSSL_VERSION_TEXT, OPENSSL_VERSION_NUMBER);

#define X(struct_name, field_name)                         \
    printf(&quot;  .&quot; #struct_name &quot;_&quot; #field_name &quot; = %ld,\n&quot;, \
           offsetof(struct struct_name, field_name));
    SSL_STRUCT_OFFSETS
#undef X
    printf(&quot;};\n&quot;);
    return 0;
}
</code></pre>

<h3>参数提取</h3>

<p>对于参数，你需要确认被HOOK函数的参数类型，以便确认读取方式，可以参考<code>kern/openssl_kern.c</code>内的<code>SSL_write</code>函数实现。</p>

<h2>编译</h2>

<h3>ARM Linux 编译</h3>

<p>公有云厂商大部分都提供了ARM64 CPU服务器，笔者选择了腾讯云的。在<code>广州六区</code>中，名字叫<code>标准型SR1</code>(SR1即ARM 64CPU)，最低配的<code>SR1.MEDIUM2</code> 2核2G即满足编译环境。可以按照<code>按量计费</code>方式购买，随时释放，比较划算。</p>

<p>操作系统选择<code>ubuntu 20.04 arm64</code>。</p>

<pre><code class="language-shell">ubuntu@VM-0-5-ubuntu:~$sudo apt-get update
ubuntu@VM-0-5-ubuntu:~$sudo apt-get install --yes wget git golang build-essential pkgconf libelf-dev llvm-12 clang-12  linux-tools-generic linux-tools-common
ubuntu@VM-0-5-ubuntu:~$wget https://golang.google.cn/dl/go1.18.linux-arm64.tar.gz
ubuntu@VM-0-5-ubuntu:~$sudo rm -rf /usr/local/go &amp;&amp; sudo tar -C /usr/local -xzf go1.18.linux-arm64.tar.gz
ubuntu@VM-0-5-ubuntu:~$for tool in &quot;clang&quot; &quot;llc&quot; &quot;llvm-strip&quot;
do
sudo rm -f /usr/bin/$tool
sudo ln -s /usr/bin/$tool-12 /usr/bin/$tool
done
ubuntu@VM-0-5-ubuntu:~$export GOPROXY=https://goproxy.cn
ubuntu@VM-0-5-ubuntu:~$export PATH=$PATH:/usr/local/go/bin
</code></pre>

<h3>编译方法</h3>

<ol>
<li><code>ANDROID=1 make</code> 命令编译支持core版本的二进制程序。</li>
<li><code>ANDROID=1 make nocore</code>命令编译仅支持当前内核版本的二进制程序。</li>
</ol>

<h1>代码仓库</h1>

<ol>
<li><a href="https://github.com/ehids/ecapture">https://github.com/ehids/ecapture</a></li>
<li><a href="https://ecapture.cc">https://ecapture.cc</a></li>
</ol>

<h2>贡献者</h2>

<p>感谢chriskaliX 、chenhengqi、vincentmli、huzai9527、yihong0618、sfx(4ft35t)等朋友的贡献。
感谢tiann Weishu在Android这个需求上的推动。
<img src="https://www.cnxct.com/wp-content/uploads/2022/09/github.com_ehids_ecapture.png" alt="" /></p>

<h1>招聘</h1>

<p>Leader直招，没有中间商赚差价。面向RASP、HIDS产品，JVM、Linux内核入侵检测等安全研发职位</p>

<ol>
<li>JAVA资深研发工程师</li>
<li>golang高级工程师</li>
<li>更多职位见 <a href="https://www.cnxct.com/jobs/?f=wxg">https://www.cnxct.com/jobs/</a></li>
</ol>
`

// https://zhuanlan.zhihu.com/p/84822157
func TestPson2(t *testing.T) {
	info, err := HtmlToSlateJson(newPTag2)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}

var newPTag21 = `
<p>现在我需要在一次事务中进行两次select</p>

<pre><code class="language-go">tx, err := mysql.Begin()
if err != nil {
	return
}

//进行第一次查询
rows, err := tx.Query(&amp;quot;select ....&amp;quot;)
if err != nil {
	return
}
defer rows.Close()

//使用第一次查询结果进行第二次查询
var result string
for rows.Next{
	if err = rows.Scan(&amp;amp;result); err != nil {
		return
	}
	//第二次Query时会报错
	rows2, err := tx.Query(&amp;quot;select ....&amp;quot;, result)
}

</code></pre>

<p>请问有什么方法可以解决？</p>
`

// https://zhuanlan.zhihu.com/p/84822157
func TestPso3(t *testing.T) {
	info, err := HtmlToSlateJson(newPTag21)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}

var errorContent = `
<p>

</p>
`

func TestErrorContent(t *testing.T) {
	info, err := HtmlToSlateJson(errorContent)
	assert.NoError(t, err)
	fmt.Printf("info--------------->"+"%+v\n", info)
}
