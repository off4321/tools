Option Explicit

' マーメイド変換ボタンをE12セルに追加
Public Sub AddMermaidConversionButton()
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Sheets("情報設定")
    
    ' 既存のボタンを確認して削除
    Dim shp As Shape
    For Each shp In ws.Shapes
        If shp.Name = "btnMermaidConversion" Then
            shp.Delete
            Exit For
        End If
    Next shp
    
    ' マーメイド変換ボタンを追加
    Dim btnMermaid As Shape
    Set btnMermaid = ws.Shapes.AddFormControl(xlButtonControl, Range("E1").Left, Range("E1").Top, Range("E1").Width, Range("E1").Height)
    btnMermaid.Name = "btnMermaidConversion"
    btnMermaid.TextFrame.Characters.Text = "マーメイド変換"
    btnMermaid.OnAction = "ConvertMermaidToExcel"
End Sub

' メイン関数：マーメイド図をExcelに変換
Sub ConvertMermaidToExcel()
    ' マークダウンファイル選択ダイアログを開く
    Dim mdFilePath As String
    mdFilePath = BrowseMdFile()
    
    If mdFilePath = "" Then
        MsgBox "ファイルが選択されませんでした", vbExclamation
        Exit Sub
    End If
    
    ' マークダウンファイルからparticipantを抽出
    Dim participants() As String
    participants = ExtractParticipants(mdFilePath)
    
    If UBound(participants) < 0 Then
        MsgBox "マーメイドシーケンス図のparticipantsが見つかりませんでした", vbExclamation
        Exit Sub
    End If
    
    ' タイムスタンプ付きの新しいシートを作成
    Dim newSheetName As String
    newSheetName = CreateNewSheet()
    
    ' メッセージとタイムスタンプを抽出
    Dim messages As Collection
    Set messages = ExtractMessagesAndTimestamps(mdFilePath)
    
    ' 最大の垂直位置を計算（垂直線の長さ用）
    Dim maxVerticalPos As Long
    maxVerticalPos = CalculateMaxVerticalPosition(messages.Count)
    
    ' 参加者を四角形として新しいシートに描画
    Dim participantPositions As Object
    Set participantPositions = DrawParticipants(participants, newSheetName, maxVerticalPos)
    
    ' 矢印を描画
    DrawArrows messages, participantPositions, newSheetName
    
    MsgBox "マーメイドシーケンス図の変換が完了しました。" & vbCrLf & _
           (UBound(participants) - LBound(participants) + 1) & "個のparticipantと" & messages.Count & "個のメッセージを変換しました。", vbInformation
End Sub

' 最大垂直位置を計算
Function CalculateMaxVerticalPosition(messageCount As Long) As Long
    Const VERTICAL_START As Long = 100
    Const VERTICAL_SPACING As Long = 30
    Const ADDITIONAL_SPACE As Long = 100 ' 追加の余白
    
    CalculateMaxVerticalPosition = VERTICAL_START + (messageCount * VERTICAL_SPACING) + ADDITIONAL_SPACE
End Function

' マークダウンファイルの参照と選択
Function BrowseMdFile() As String
    Dim fd As FileDialog
    Set fd = Application.FileDialog(msoFileDialogFilePicker)
    
    With fd
        .Title = "マーメイドシーケンス図のMarkdownファイルを選択"
        .AllowMultiSelect = False
        .Filters.Clear
        .Filters.Add "Markdown Files", "*.md"
        If .Show = -1 Then
            BrowseMdFile = .SelectedItems(1)
        Else
            BrowseMdFile = ""
        End If
    End With
    
    Set fd = Nothing
End Function

' マーメイドシーケンス図からparticipantを抽出
Function ExtractParticipants(filePath As String) As String()
    Dim fso As Object
    Dim textFile As Object
    Dim content As String
    Dim lines() As String
    Dim i As Long
    Dim inMermaidBlock As Boolean
    Dim participantList As New Collection
    
    ' ファイルシステムオブジェクトを作成
    Set fso = CreateObject("Scripting.FileSystemObject")
    
    ' ファイルの存在確認
    If Not fso.FileExists(filePath) Then
        MsgBox "ファイルが見つかりません: " & filePath, vbExclamation
        Dim emptyArray() As String
        ExtractParticipants = emptyArray
        Exit Function
    End If
    
    ' ファイルを開いてコンテンツを読み込む
    Set textFile = fso.OpenTextFile(filePath, 1) ' 1 = ForReading
    content = textFile.ReadAll
    textFile.Close
    
    ' コンテンツを行に分割
    lines = Split(content, vbLf)
    
    ' マーメイドシーケンス図内のparticipant行を検索
    inMermaidBlock = False
    For i = LBound(lines) To UBound(lines)
        Dim line As String
        line = Trim(lines(i))
        
        ' マーメイドブロックの開始をチェック
        If line = "```mermaid" Then
            inMermaidBlock = True
        ElseIf line = "```" And inMermaidBlock Then
            inMermaidBlock = False
        ElseIf inMermaidBlock And Left(Trim(line), 11) = "participant " Then
            ' participant名を抽出 - スペース後の全テキスト
            Dim participant As String
            participant = Trim(Mid(line, 12))
            
            ' コレクションに追加（重複なし）
            On Error Resume Next
            participantList.Add participant, participant
            On Error GoTo 0
        ElseIf inMermaidBlock And (InStr(line, "->") > 0 Or InStr(line, "->>") > 0) Then
            ' 矢印の種類を検出して正確に分割
            Dim arrowType As String
            If InStr(line, "->>") > 0 Then
                arrowType = "->>"
            Else
                arrowType = "->"
            End If
            
            ' メッセージ行（sender->receiver）から抽出
            Dim parts() As String
            parts = Split(line, arrowType)
            If UBound(parts) >= 0 Then
                Dim sender As String, receiver As String
                sender = Trim(parts(0))
                
                ' 送信者をコレクションに追加（重複なし）
                On Error Resume Next
                participantList.Add sender, sender
                On Error GoTo 0
                
                ' 受信者を抽出
                If UBound(parts) > 0 Then
                    receiver = Trim(parts(1))
                    ' 最初のコロンより後のテキストを削除
                    If InStr(receiver, ":") > 0 Then
                        receiver = Trim(Left(receiver, InStr(receiver, ":") - 1))
                    End If
                    
                    ' 受信者をコレクションに追加（重複なし）
                    On Error Resume Next
                    participantList.Add receiver, receiver
                    On Error GoTo 0
                End If
            End If
        End If
    Next i
    
    ' コレクションを配列に変換
    Dim resultArray() As String
    If participantList.Count > 0 Then
        ReDim resultArray(0 To participantList.Count - 1)
        For i = 1 To participantList.Count
            resultArray(i - 1) = participantList(i)
        Next i
    Else
        ReDim resultArray(-1 To -1)
    End If
    
    ExtractParticipants = resultArray
End Function

' タイムスタンプ付きの新しいシートを作成
Function CreateNewSheet() As String
    Dim newSheetName As String
    newSheetName = "シーケンス_" & Format(Now, "yyyymmdd_hhmmss")
    
    ' 新しいワークシートを追加
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Sheets.Add(After:=ThisWorkbook.Sheets(ThisWorkbook.Sheets.Count))
    ws.Name = newSheetName
    
    CreateNewSheet = newSheetName
End Function

' participantを四角形として描画し、位置情報を返す
Function DrawParticipants(participants() As String, sheetName As String, Optional verticalLength As Long = 600) As Object
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Sheets(sheetName)
    
    ' 初期位置とサイズの設定
    Const TOP_MARGIN As Long = 20
    Const LEFT_MARGIN As Long = 50
    Const BOX_WIDTH As Long = 100
    Const BOX_HEIGHT As Long = 50
    Const SPACING As Long = 150
    
    ' 垂直線の長さ - 動的に調整
    Dim LINE_LENGTH As Long
    LINE_LENGTH = verticalLength  ' 渡された長さまたはデフォルトの600
    
    ' participant位置情報を格納する Dictionary
    Dim positions As Object
    Set positions = CreateObject("Scripting.Dictionary")
    
    ' まずタイムラインの四角形を一番左に描画
    Dim timelineShape As Shape
    Set timelineShape = ws.Shapes.AddShape(msoShapeRectangle, LEFT_MARGIN, TOP_MARGIN, BOX_WIDTH, BOX_HEIGHT)
    
    ' タイムラインの四角形の設定
    With timelineShape
        .Name = "Timeline"
        .Fill.ForeColor.RGB = RGB(255, 220, 200)  ' 薄いオレンジ色
        .Line.ForeColor.RGB = RGB(0, 0, 0)
        .TextFrame.Characters.Text = "Timeline"
        .TextFrame.HorizontalAlignment = xlHAlignCenter
        .TextFrame.VerticalAlignment = xlVAlignCenter
        .TextFrame.Characters.Font.Bold = True
    End With
    
    ' タイムライン位置情報を登録
    Dim timelinePos As Object
    Set timelinePos = CreateObject("Scripting.Dictionary")
    timelinePos.Add "left", LEFT_MARGIN
    timelinePos.Add "center", LEFT_MARGIN + (BOX_WIDTH / 2)
    timelinePos.Add "top", TOP_MARGIN
    timelinePos.Add "bottom", TOP_MARGIN + BOX_HEIGHT
    positions.Add "Timeline", timelinePos
    
    ' タイムラインから垂直線を引く
    Dim timelineLine As Shape
    Set timelineLine = ws.Shapes.AddLine(LEFT_MARGIN + (BOX_WIDTH / 2), TOP_MARGIN + BOX_HEIGHT, _
                                       LEFT_MARGIN + (BOX_WIDTH / 2), TOP_MARGIN + BOX_HEIGHT + LINE_LENGTH)
    With timelineLine
        .Name = "TimelineLine"
        .Line.ForeColor.RGB = RGB(0, 0, 0)
        .Line.Weight = 1.5
        .Line.DashStyle = msoLineDash  ' 破線スタイル
    End With
    
    ' 各participantに対して四角形を描画（タイムラインの右側に配置）
    Dim i As Long
    Dim left As Long
    
    For i = LBound(participants) To UBound(participants)
        ' タイムラインの右側から配置するため、位置を調整
        left = LEFT_MARGIN + BOX_WIDTH + SPACING + (i * SPACING)
        
        ' シェイプの追加
        Dim shp As Shape
        Set shp = ws.Shapes.AddShape(msoShapeRectangle, left, TOP_MARGIN, BOX_WIDTH, BOX_HEIGHT)
        
        ' シェイプの設定
        With shp
            .Name = "Participant_" & i
            .Fill.ForeColor.RGB = RGB(200, 200, 255)
            .Line.ForeColor.RGB = RGB(0, 0, 0)
            .TextFrame.Characters.Text = participants(i)
            .TextFrame.HorizontalAlignment = xlHAlignCenter
            .TextFrame.VerticalAlignment = xlVAlignCenter
        End With
        
        ' participant位置情報を登録
        Dim participantPos As Object
        Set participantPos = CreateObject("Scripting.Dictionary")
        participantPos.Add "left", left
        participantPos.Add "center", left + (BOX_WIDTH / 2)
        participantPos.Add "top", TOP_MARGIN
        participantPos.Add "bottom", TOP_MARGIN + BOX_HEIGHT
        positions.Add participants(i), participantPos
        
        ' 各参加者の四角形から垂直線を引く
        Dim participantLine As Shape
        Set participantLine = ws.Shapes.AddLine(left + (BOX_WIDTH / 2), TOP_MARGIN + BOX_HEIGHT, _
                                             left + (BOX_WIDTH / 2), TOP_MARGIN + BOX_HEIGHT + LINE_LENGTH)
        With participantLine
            .Name = "ParticipantLine_" & i
            .Line.ForeColor.RGB = RGB(0, 0, 0)
            .Line.Weight = 1.5
            .Line.DashStyle = msoLineDashDot  ' 破線スタイル
        End With
    Next i
    
    Set DrawParticipants = positions
End Function

' マーメイドシーケンス図からメッセージとタイムスタンプを抽出
Function ExtractMessagesAndTimestamps(filePath As String) As Collection
    Dim fso As Object
    Dim textFile As Object
    Dim content As String
    Dim lines() As String
    Dim i As Long
    Dim inMermaidBlock As Boolean
    Dim messages As New Collection
    Dim currentMessage As Object  ' Dictionary型をObject型に変更
    Dim currentTimestamp As String
    
    ' ファイルシステムオブジェクトを作成
    Set fso = CreateObject("Scripting.FileSystemObject")
    
    ' ファイルを開いてコンテンツを読み込む
    Set textFile = fso.OpenTextFile(filePath, 1) ' 1 = ForReading
    content = textFile.ReadAll
    textFile.Close
    
    ' コンテンツを行に分割
    lines = Split(content, vbLf)
    
    ' マーメイドシーケンス図内のメッセージ行とタイムスタンプを検索
    inMermaidBlock = False
    For i = LBound(lines) To UBound(lines)
        Dim line As String
        line = Trim(lines(i))
        
        ' マーメイドブロックの開始をチェック
        If line = "```mermaid" Then
            inMermaidBlock = True
        ElseIf line = "```" And inMermaidBlock Then
            inMermaidBlock = False
        ElseIf inMermaidBlock Then
            ' メッセージ行（sender->receiver）を検出
            If InStr(line, "->") > 0 Or InStr(line, "->>") > 0 Then
                Set currentMessage = CreateObject("Scripting.Dictionary")  ' New Dictionary から変更
                
                ' 矢印の種類を検出して正確に分割
                Dim arrowType As String
                If InStr(line, "->>") > 0 Then
                    arrowType = "->>"
                    currentMessage.Add "arrow_type", "async"
                Else
                    arrowType = "->"
                    currentMessage.Add "arrow_type", "sync"
                End If
                
                ' メッセージ行を分解して送信元・送信先・メッセージ内容を抽出
                Dim parts() As String
                parts = Split(line, arrowType)
                
                If UBound(parts) >= 0 Then
                    Dim sender As String, receiver As String, msg As String
                    sender = Trim(parts(0))
                    currentMessage.Add "sender", sender
                    
                    If UBound(parts) > 0 Then
                        ' 受信者とメッセージを分離
                        Dim receiverParts() As String
                        receiverParts = Split(parts(1), ":", 2)
                        
                        receiver = Trim(receiverParts(0))
                        currentMessage.Add "receiver", receiver
                        
                        ' メッセージが存在すれば抽出
                        If UBound(receiverParts) > 0 Then
                            msg = Trim(receiverParts(1))
                        Else
                            msg = ""
                        End If
                        currentMessage.Add "message", msg
                    End If
                End If
                
                ' タイムスタンプはまだ設定しない
                currentMessage.Add "timestamp", ""
                
                ' コレクションに追加
                messages.Add currentMessage
            
            ' タイムスタンプ行の検出
            ElseIf Left(Trim(line), 12) = "Note over Ti" And InStr(line, "Timeline:") > 0 Then
                ' タイムスタンプを抽出
                currentTimestamp = Trim(Mid(line, InStr(line, ":") + 1))
                
                ' 日付と時刻部分のみを抽出する（日付時刻パターンに基づいて）
                Dim timestampParts As Variant
                Dim cleanedTimestamp As String
                
                ' 日付/時刻パターンを認識 - 最初の数値と文字列の部分のみを取得
                ' 例：Mar 23, 2025 12:08:35.701019000 東京 (標準時) -> Mar 23, 2025 12:08:35.701019000
                
                ' 正規表現が使えないので空白で分割して最初の数個の要素だけを使用
                timestampParts = Split(currentTimestamp, " ")
                
                ' タイムスタンプの典型的なパターン: [月] [日], [年] [時:分:秒.ミリ秒]
                ' 最大4つのパーツ（月、日、年、時間）までを保持し、それ以降は無視
                Dim maxParts As Integer
                maxParts = 4 ' 月、日、年、時間部分まで
                
                If UBound(timestampParts) >= 0 Then
                    cleanedTimestamp = timestampParts(0) ' 月
                    
                    ' 2番目のパーツが日付（数字+カンマ）の場合
                    If UBound(timestampParts) >= 1 Then
                        cleanedTimestamp = cleanedTimestamp & " " & timestampParts(1) ' 日,
                        
                        ' 3番目のパーツが年
                        If UBound(timestampParts) >= 2 Then
                            cleanedTimestamp = cleanedTimestamp & " " & timestampParts(2) ' 年
                            
                            ' 4番目のパーツが時間
                            If UBound(timestampParts) >= 3 And IsTimeFormat(timestampParts(3)) Then
                                cleanedTimestamp = cleanedTimestamp & " " & timestampParts(3) ' 時間
                                
                                ' ミリ秒部分がある場合（5番目のパーツが数値のみで構成されている）
                                If UBound(timestampParts) >= 4 And IsNumeric(Replace(timestampParts(4), ".", "")) Then
                                    cleanedTimestamp = cleanedTimestamp & " " & timestampParts(4) ' ミリ秒
                                End If
                            End If
                        End If
                    End If
                Else
                    cleanedTimestamp = currentTimestamp ' 分割できない場合は元の値を使用
                End If
                
                ' 直前のメッセージにタイムスタンプを設定
                If messages.Count > 0 Then
                    Set currentMessage = messages(messages.Count)
                    currentMessage("timestamp") = cleanedTimestamp
                End If
            End If
        End If
    Next i
    
    Set ExtractMessagesAndTimestamps = messages
End Function

' 文字列が時間形式（HH:MM:SS）かどうかをチェックする関数
Function IsTimeFormat(ByVal strTime As String) As Boolean
    ' 時間形式（数字:数字:数字）かどうかをチェック
    Dim parts As Variant
    parts = Split(strTime, ":")
    
    ' 最低でも時間と分（HH:MM）があるはず
    If UBound(parts) < 1 Then
        IsTimeFormat = False
        Exit Function
    End If
    
    ' 各部分が数値であることを確認
    Dim i As Integer
    For i = 0 To UBound(parts)
        ' 最後の部分は秒とミリ秒かもしれないので、小数点を許容
        If i = UBound(parts) Then
            If Not IsNumeric(Replace(parts(i), ".", "")) Then
                IsTimeFormat = False
                Exit Function
            End If
        ElseIf Not IsNumeric(parts(i)) Then
            IsTimeFormat = False
            Exit Function
        End If
    Next i
    
    IsTimeFormat = True
End Function

' 矢印を描画する
Sub DrawArrows(messages As Collection, positions As Object, sheetName As String)
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Sheets(sheetName)
    
    ' 垂直位置のスタート地点と間隔
    Const VERTICAL_START As Long = 100
    Const VERTICAL_SPACING As Long = 30
    
    ' タイムスタンプとメッセージテキストの位置調整用
    Const TIMESTAMP_OFFSET_X As Long = 10
    Const TIMESTAMP_OFFSET_Y As Long = -5
    Const MESSAGE_OFFSET_Y_SINGLE As Long = -25  ' 1行の場合の上部オフセット
    Const MESSAGE_OFFSET_Y_MULTI As Long = -35   ' 複数行の場合の上部オフセット
    
    ' メッセージテキストボックスのサイズ設定
    Const MESSAGE_WIDTH As Long = 200  ' より広いテキストボックス
    Const MESSAGE_HEIGHT As Long = 40  ' より高いテキストボックス
    Const TIMESTAMP_WIDTH As Long = 300 ' タイムスタンプ用テキストボックス幅（拡大）
    
    Dim i As Long
    Dim currentVerticalPos As Long
    currentVerticalPos = VERTICAL_START
    
    ' 各メッセージに対して矢印を描画
    For i = 1 To messages.Count
        Dim msg As Object
        Set msg = messages(i)
        
        Dim sender As String, receiver As String, message As String, timestamp As String
        sender = msg("sender")
        receiver = msg("receiver")
        message = msg("message")
        timestamp = msg("timestamp")
        
        ' 送信元と受信先が位置情報に存在する場合のみ処理
        If positions.Exists(sender) And positions.Exists(receiver) Then
            ' 送信元と受信先の位置を取得
            Dim senderPos As Object, receiverPos As Object
            Set senderPos = positions(sender)
            Set receiverPos = positions(receiver)
            
            ' 矢印の開始点と終了点
            Dim startX As Long, startY As Long, endX As Long, endY As Long
            startX = senderPos("center")
            startY = currentVerticalPos
            endX = receiverPos("center")
            endY = currentVerticalPos
            
            ' 矢印を描画
            Dim arrow As Shape
            Set arrow = ws.Shapes.AddLine(startX, startY, endX, endY)
            
            ' 矢印のスタイル設定
            With arrow
                .Name = "Arrow_" & i
                .Line.ForeColor.RGB = RGB(0, 0, 150)
                .Line.Weight = 1.5
                
                ' 矢印のタイプ設定
                If msg("arrow_type") = "async" Then
                    .Line.EndArrowheadStyle = msoArrowheadOpen
                Else
                    .Line.EndArrowheadStyle = msoArrowheadTriangle
                End If
                .Line.EndArrowheadLength = msoArrowheadLengthMedium
                .Line.EndArrowheadWidth = msoArrowheadWidthMedium
            End With
            
            ' メッセージに改行が含まれているかチェック
            Dim messageOffsetY As Long
            If InStr(message, "<br>") > 0 Or InStr(message, vbCrLf) > 0 Then
                messageOffsetY = MESSAGE_OFFSET_Y_MULTI  ' 複数行の場合はより上に配置
            Else
                messageOffsetY = MESSAGE_OFFSET_Y_SINGLE ' 1行の場合は標準位置
            End If
            
            ' メッセージテキストを追加（矢印の上に配置）
            Dim messageText As Shape
            Set messageText = ws.Shapes.AddTextbox(msoTextOrientationHorizontal, _
                                                 (startX + endX) / 2 - (MESSAGE_WIDTH / 2), _
                                                 currentVerticalPos + messageOffsetY, _
                                                 MESSAGE_WIDTH, MESSAGE_HEIGHT)
            With messageText
                .Name = "MessageText_" & i
                ' <br>タグを改行に置換
                message = Replace(message, "<br>", vbCrLf)
                .TextFrame.Characters.Text = message
                .TextFrame.HorizontalAlignment = xlHAlignCenter
                .TextFrame.VerticalAlignment = xlVAlignCenter
                .Fill.Visible = msoFalse
                .Line.Visible = msoFalse
                .TextFrame.Characters.Font.Size = 8
                ' 自動サイズ調整を有効化
                .TextFrame.AutoSize = True
            End With
            
            ' タイムスタンプがあれば、Timeline上に表示
            If timestamp <> "" Then
                ' タイムスタンプ用テキストボックスを追加
                Dim timestampBox As Shape
                Set timestampBox = ws.Shapes.AddTextbox(msoTextOrientationHorizontal, _
                                                     positions("Timeline")("left") + TIMESTAMP_OFFSET_X, _
                                                     currentVerticalPos + TIMESTAMP_OFFSET_Y, _
                                                     TIMESTAMP_WIDTH, 20)
                With timestampBox
                    .Name = "Timestamp_" & i
                    .TextFrame.Characters.Text = timestamp
                    .TextFrame.HorizontalAlignment = xlHAlignLeft
                    .TextFrame.VerticalAlignment = xlVAlignCenter
                    .Fill.Visible = msoFalse
                    .Line.Visible = msoFalse
                    .TextFrame.Characters.Font.Size = 8
                    .TextFrame.Characters.Font.Italic = True
                End With
            End If
            
            ' 次のメッセージの垂直位置を更新
            currentVerticalPos = currentVerticalPos + VERTICAL_SPACING
        End If
    Next i
End Sub
