# Redeeman

Redeeman is a program written in Go that automates the process of logging into a specific website and redeeming codes from txt files. It interacts with the target website by using Selenium WebDriver and ChromeDriver.

## Prerequisites

The following components are required:

1. **Go** - The Go programming language, which you can download from [here](https://golang.org/dl/).
2. **Selenium WebDriver** - You can download selenium-server-standalone-3.141.59.jar from [here](https://www.selenium.dev/downloads/).
3. **ChromeDriver** - You can download it from [here](https://sites.google.com/chromium.org/driver/).

## Usage

Clone this repository, navigate to the directory, and run the following command to start the program:

```bash
go run main.go
```

Make sure to populate the .env file with your credentials:

```env
ACCESS_POKEMON_ID=your_id_here
ACCESS_POKEMON_PASSWORD=your_password_here
```

Place the txt files containing the codes you want to redeem in the targets directory.

## Important Notes
1. **Website Changes** - The program may stop working correctly if the layout of the website changes. It depends on specific element IDs and XPaths that may change over time.
2. **Rate Limiting and Invalid Codes** - The program pauses for a few seconds after each attempt to redeem a code. If a code is invalid, or if too many attempts are made in a short period of time, the program may not work as expected.
3. **reCAPTCHA Validation** - If the message "The server was not able to validate your reCAPTCHA submission." appears, the program might not work correctly.
4. **Headless Mode** - The headless mode might not function correctly. For this reason, it is commented out in the code.
5. **Terms of Service** - Be careful when using the program, as it should not violate the terms of service of the website or any other related policies.
## Contributions
Contributions are welcome. Please open an issue if you find a bug or have a feature request.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

<br><br>

---------------------------------------------------------------
<br>

# Redeeman

Redeemanは、特定のウェブサイトにログインし、txtファイルからコードを自動的に登録するプロセスを自動化するGo言語で書かれたプログラムです。Selenium WebDriverとChromeDriverを使用して対象のウェブサイトと対話します。

## 必要なコンポーネント

次のコンポーネントが必要です：

1. **Go** - Goプログラミング言語。[ここ](https://golang.org/dl/)からダウンロードできます。
2. **Selenium WebDriver** - selenium-server-standalone-3.141.59.jarは[ここ](https://www.selenium.dev/downloads/)からダウンロードできます。
3. **ChromeDriver** - [ここ](https://sites.google.com/chromium.org/driver/)からダウンロードできます。

## 使用法

このリポジトリをクローンし、ディレクトリに移動した後、次のコマンドを実行してプログラムを開始します：

```bash
go run main.go
```

.env ファイルに認証情報を記入してください：
```env
ACCESS_POKEMON_ID=あなたのID
ACCESS_POKEMON_PASSWORD=あなたのパスワード
```

登録したいコードを含むtxtファイルを targets ディレクトリに配置します。

## 重要な注意事項
1. **ウェブサイトの変更** - ウェブサイトのレイアウトが変更されると、プログラムは正しく動作しなくなる可能性があります。これは、時間の経過と共に変更される可能性のある特定の要素IDとXPathsに依存しています。
2. **レート制限と無効なコード** - プログラムは、コードの登録を試みるたびに数秒間一時停止します。コードが無効であったり、短時間に多くの試行が行われた場合、プログラムは予期したように動作しないかもしれません。
3. **reCAPTCHAの検証** - "The server was not able to validate your reCAPTCHA submission."というメッセージが表示されると、プログラムは正常に動作しない可能性があります。
4. **ヘッドレスモード** - ヘッドレスモードは正常に動作しないかもしれません。そのため、コード内ではコメントアウトされています。
5. **利用規約** - プログラムの使用には注意が必要です。ウェブサイトの利用規約やその他の関連ポリシーを違反しないようにしてください。
## 貢献
バグが見つかった場合や、機能リクエストがある場合は、Issueを開いてください。

## ライセンス
このプロジェクトはMITライセンスの下でライセンスされています。詳細はLICENSEファイルをご覧ください。
