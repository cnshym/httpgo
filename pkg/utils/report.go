package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// URLFingerprint 结构体表示每个 URL 的指纹信息
type URLFingerprint struct {
	Url        string
	StatusCode int
	Title      string
	CmsList    string
	OtherList  string
	Screenshot string
}

var HtmlHeaderA = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>httpgo Fingerprint Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
            color: #333;
        }
        h1 {
            text-align: center;
            margin: 20px 0;
            color: #444;
        }
        table {
            width: 90%;
            margin: 20px auto;
            border-collapse: collapse;
            background: #fff;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        table, th, td {
            border: 1px solid #ddd;
        }
        th, td {
            padding: 10px;
            text-align: left;
        }
        th {
            background-color: #f8f8f8;
            color: #555;
            font-size: 1rem;
        }
        .container {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            padding: 10px;
            gap: 10px;
        }
        .target {
            flex: 1;
            background: #fafafa;
            padding: 10px;
            border-radius: 8px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
            max-width: 48%;
            display: flex;
            align-items: flex-start;
            gap: 10px;
        }
        .target-info {
            flex: 1;
            max-width: calc(100% - 170px); /* 留出图片宽度和间距 */
        }
        .target-info p {
            margin: 5px 0;
            font-size: 0.9rem;
            white-space: normal;
            word-wrap: break-word;
        }
        .target-img {
            flex: 0 0 150px;
            text-align: center;
        }
        .target-img img {
            max-width: 150px;
            height: auto;
            border-radius: 8px;
            cursor: pointer;
            transition: opacity 0.3s;
        }
        .target-img img:hover {
            opacity: 0.8;
        }
        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.8);
            align-items: center;
            justify-content: center;
            z-index: 1000;
        }
        .modal-content {
            max-width: 90%;
            max-height: 90%;
            position: relative;
        }
        .modal-content img {
            width: 100%;
            height: auto;
            border: 5px solid #fff;
            border-radius: 8px;
        }
        .modal-close {
            position: absolute;
            top: 20px;
            right: 20px;
            font-size: 2rem;
            color: #fff;
            cursor: pointer;
            transition: color 0.3s;
        }
        .modal-close:hover {
            color: #ddd;
        }
        .cms-info {
            color: red;
        }
        .other-info {
            color: green;
        }
        .stats {
            margin: 20px auto;
            width: 90%;
            padding: 15px;
            background: #fafafa;
            border-radius: 8px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
        .stats h2 {
            margin-top: 0;
            font-size: 1.2rem;
        }
        .button-group {
            display: flex;
            flex-wrap: wrap;
            margin: 10px 0;
        }
        .button-group button {
            background-color: #007bff;
            color: white;
            border: none;
            padding: 6px 12px;
            margin: 4px;
            border-radius: 4px;
            cursor: pointer;
            transition: background-color 0.3s;
            font-size: 0.875rem;
        }
        .button-group button:hover {
            background-color: #0056b3;
        }
        #scroll-to-top, #copy-urls {
            position: fixed;
            bottom: 20px;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            padding: 10px 15px;
            cursor: pointer;
            font-size: 0.875rem;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
            transition: background-color 0.3s, box-shadow 0.3s;
        }
        #scroll-to-top {
            right: 20px;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 18px;
            display: none;
        }
        #copy-urls {
            right: 70px;
        }
        #scroll-to-top:hover, #copy-urls:hover {
            background-color: #0056b3;
            box-shadow: 0 6px 12px rgba(0, 0, 0, 0.3);
        }
    </style>
    <script>
        document.addEventListener("DOMContentLoaded", function() {
            const scrollToTopButton = document.getElementById("scroll-to-top");
            const copyUrlsButton = document.getElementById("copy-urls");
            const modal = document.getElementById("modal");
            const modalImg = document.getElementById("modal-img");
            const modalClose = document.querySelector(".modal-close");
            let originalData = [];
            let currentData = [];

            // Scroll to top
            scrollToTopButton.addEventListener("click", function() {
                window.scrollTo({ top: 0, behavior: "smooth" });
            });

            window.addEventListener("scroll", function() {
                scrollToTopButton.style.display = window.scrollY > 300 ? "flex" : "none";
            });

            // Modal functions
            function openModal(src) {
                modal.style.display = "flex";
                modalImg.src = src;
            }

            function closeModal() {
                modal.style.display = "none";
            }

            modalClose.addEventListener("click", closeModal);
            modal.addEventListener("click", function(event) {
                if (event.target === modal) {
                    closeModal();
                }
            });

            document.addEventListener("click", function(event) {
                if (event.target.tagName === "IMG" && event.target.closest(".target-img")) {
                    openModal(event.target.src);
                }
            });

            // Copy URLs to clipboard
            function copyURLs() {
                if (currentData.length === 0) {
                    alert("没有可复制的 URL！");
                    return;
                }
                const urls = currentData.map(item => item.Url).join("\n");
                // 优先使用 navigator.clipboard
                if (navigator.clipboard && navigator.clipboard.writeText) {
                    navigator.clipboard.writeText(urls).then(() => {
                        const originalText = copyUrlsButton.textContent;
                        copyUrlsButton.textContent = "已复制!";
                        setTimeout(() => {
                            copyUrlsButton.textContent = originalText;
                        }, 2000);
                    }).catch(err => {
                        console.error("Clipboard API 复制失败:", err);
                        fallbackCopy(urls);
                    });
                } else {
                    fallbackCopy(urls);
                }
            }

            // 备用复制方法
            function fallbackCopy(text) {
                try {
                    const textarea = document.createElement("textarea");
                    textarea.value = text;
                    document.body.appendChild(textarea);
                    textarea.select();
                    document.execCommand("copy");
                    document.body.removeChild(textarea);
                    const originalText = copyUrlsButton.textContent;
                    copyUrlsButton.textContent = "已复制!";
                    setTimeout(() => {
                        copyUrlsButton.textContent = originalText;
                    }, 2000);
                } catch (err) {
                    console.error("备用复制失败:", err);
                    alert("复制失败，请检查浏览器权限或手动复制！");
                }
            }

            copyUrlsButton.addEventListener("click", copyURLs);

            // Update statistics
            function updateStats(data) {
                const cmsCount = {};
                const otherCount = {};
                const statusCodeCount = {};

                data.forEach(item => {
                    item.CmsList.split(";").forEach(cms => {
                        cms = cms.trim();
                        if (cms) cmsCount[cms] = (cmsCount[cms] || 0) + 1;
                    });
                    item.OtherList.split(";").forEach(other => {
                        other = other.trim();
                        if (other) otherCount[other] = (otherCount[other] || 0) + 1;
                    });
                    const statusCode = item.StatusCode;
                    if (statusCode) statusCodeCount[statusCode] = (statusCodeCount[statusCode] || 0) + 1;
                });

                const cmsStats = Object.entries(cmsCount).sort((a, b) => b[1] - a[1])
                    .map(([key, value]) => "<button class=\"cms-item\" data-type=\"cms\" data-value=\"" + key + "\">" + key + ": " + value + "</button>")
                    .join("");
                document.getElementById("cms-stats").innerHTML = "<h2>CMS Fingerprint Information</h2><div class=\"button-group\">" + cmsStats + "</div>";

                const otherStats = Object.entries(otherCount).sort((a, b) => b[1] - a[1])
                    .map(([key, value]) => "<button class=\"other-item\" data-type=\"other\" data-value=\"" + key + "\">" + key + ": " + value + "</button>")
                    .join("");
                document.getElementById("other-stats").innerHTML = "<br><h2>Other Fingerprint Information</h2><div class=\"button-group\">" + otherStats + "</div>";

                const statusCodeStats = Object.entries(statusCodeCount).sort((a, b) => b[1] - a[1])
                    .map(([key, value]) => "<button class=\"status-code-item\" data-type=\"status-code\" data-value=\"" + key + "\">" + key + ": " + value + "</button>")
                    .join("");
                document.getElementById("status-code-stats").innerHTML = "<br><h2>Status Code Information</h2><div class=\"button-group\">" + statusCodeStats + "</div>";

                const allCount = data.length;
                document.getElementById("all-stats").innerHTML = "<br><h2>All Fingerprint Information</h2><div class=\"button-group\"><button id=\"btn-all\">ALL (" + allCount + ")</button></div>";
            }

            // Filter data
            function filterData(data, type, value) {
                return data.filter(item => {
                    if (type === "cms") return item.CmsList.split(";").map(cms => cms.trim()).includes(value);
                    if (type === "other") return item.OtherList.split(";").map(other => other.trim()).includes(value);
                    if (type === "status-code") return item.StatusCode.toString() === value;
                    return false;
                });
            }

            // Update table
            function updateTable(data) {
                currentData = data; // 跟踪当前显示的数据
                const tableBody = document.querySelector("tbody");
                tableBody.innerHTML = "";
                for (let i = 0; i < data.length; i += 2) {
                    const leftItem = data[i];
                    const rightItem = data[i + 1] || {};
                    const row = document.createElement("tr");
                    row.innerHTML = "<td class=\"container\">" +
                        "<div class=\"target\">" +
                        (leftItem ? "<div class=\"target-info\">" +
                        "<p><strong>目标:</strong> <a href=\"" + leftItem.Url + "\" target=\"_blank\">" + leftItem.Url + "</a></p>" +
                        "<p><strong>状态码:</strong> " + leftItem.StatusCode + "</p>" +
                        "<p><strong>标题:</strong> " + leftItem.Title + "</p>" +
                        "<p><strong>CMS指纹信息:</strong> <span class=\"cms-info\">" + leftItem.CmsList + "</span></p>" +
                        "<p><strong>OTHER信息:</strong> <span class=\"other-info\">" + leftItem.OtherList + "</span></p>" +
                        "</div><div class=\"target-img\">" +
                        (leftItem.Screenshot ? "<img src=\"" + leftItem.Screenshot + "\" alt=\"Screenshot\" loading=\"lazy\">" : "<p>无截图</p>") +
                        "</div>" : "<p>没有更多数据了</p>") +
                        "</div>" +
                        "<div class=\"target\">" +
                        (rightItem.Url ? "<div class=\"target-info\">" +
                        "<p><strong>目标:</strong> <a href=\"" + rightItem.Url + "\" target=\"_blank\">" + rightItem.Url + "</a></p>" +
                        "<p><strong>状态码:</strong> " + (rightItem.StatusCode || "") + "</p>" +
                        "<p><strong>标题:</strong> " + (rightItem.Title || "") + "</p>" +
                        "<p><strong>CMS指纹信息:</strong> <span class=\"cms-info\">" + (rightItem.CmsList || "") + "</span></p>" +
                        "<p><strong>OTHER信息:</strong> <span class=\"other-info\">" + (rightItem.OtherList || "") + "</span></p>" +
                        "</div><div class=\"target-img\">" +
                        (rightItem.Screenshot ? "<img src=\"" + rightItem.Screenshot + "\" alt=\"Screenshot\" loading=\"lazy\">" : "<p>无截图</p>") +
                        "</div>" : "<p>没有更多数据了</p>") +
                        "</div></td>";
                    tableBody.appendChild(row);
                }
            }

            // Filter button click handler
            document.addEventListener("click", function(event) {
                if (event.target.classList.contains("cms-item") || event.target.classList.contains("other-item") || event.target.classList.contains("status-code-item")) {
                    const type = event.target.getAttribute("data-type");
                    const value = event.target.getAttribute("data-value");
                    updateTable(filterData(originalData, type, value));
                } else if (event.target.id === "btn-all") {
                    updateTable(originalData);
                }
            });

            // Fetch JSON data
            fetch("`

var HtmlHeaderB = `")
                .then(response => {
                    if (!response.ok) throw new Error("Network response was not ok");
                    return response.json();
                })
                .then(data => {
                    originalData = data;
                    updateStats(data);
                    updateTable(data);
                })
                .catch(error => console.error("Error loading JSON data:", error));
        });
    </script>
</head>
<body>
    <h1>URL Fingerprint Report</h1>
    <div class="stats">
        <div id="cms-stats"></div>
        <div id="other-stats"></div>
        <div id="status-code-stats"></div>
        <div id="all-stats"></div>
    </div>
    <div id="modal" class="modal">
        <div class="modal-content">
            <span class="modal-close">×</span>
            <img id="modal-img" src="" alt="Screenshot">
        </div>
    </div>
    <table>
        <thead>
            <tr>
                <th>Details</th>
            </tr>
        </thead>
        <tbody>
            <!-- Data rows will be inserted here by JavaScript -->
        </tbody>
    </table>
    <button id="copy-urls" title="Copy URLs">复制当前指纹所有URL</button>
    <button id="scroll-to-top" title="Go to Top">⇧</button>
</body>
</html>
`

// 创建 HTML 报告
func InitializeHTMLReport(filename string, json string) (*os.File, error) {
	// 拼接 HTML 头部和 JSON 文件路径
	var HtmlHeader = HtmlHeaderA + json + HtmlHeaderB
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("无法创建 HTML 文件: %v", err)
	}
	_, err = file.WriteString(HtmlHeader)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("无法写入 HTML 文件: %v", err)
	}
	return file, nil
}

// 定义全局互斥锁
var jsonMutex sync.Mutex

// AppendJSONReport 将 URLFingerprint 数据追加到指定的 JSON 文件中
func AppendJSONReport(filename string, data URLFingerprint) error {
	// 使用互斥锁确保线程安全
	jsonMutex.Lock()
	defer jsonMutex.Unlock()

	var existingData []URLFingerprint

	// 读取现有文件内容
	fileContent, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("无法读取 JSON 文件: %v", err)
	}

	// 如果文件存在且有内容，解码现有数据
	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &existingData); err != nil {
			return fmt.Errorf("无法解码 JSON 内容: %v", err)
		}
	}

	// 添加新数据
	existingData = append(existingData, data)

	// 打开文件进行写入（覆盖模式）
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("无法打开 JSON 文件: %v", err)
	}
	defer file.Close()

	// 创建 JSON 编码器并设置缩进
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	// 写入数据
	if err := encoder.Encode(existingData); err != nil {
		return fmt.Errorf("无法写入 JSON 数据: %v", err)
	}

	return nil
}
