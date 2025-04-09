Option Explicit


Sub GenerateMermaidSequence()
    Dim cmd As String
    Dim exePath As String
    
    exePath = ThisWorkbook.Path & "\..\packet_sequence.exe"
    If Dir(exePath) = "" Then
        MsgBox "packet_sequence.exe not found in the parent directory.", vbExclamation
        Exit Sub
    End If
    
    Dim filePath As String: filePath = Range("D3").Value
    Dim outPath As String: outPath = Range("D4").Value
    Dim maxPkts As String: maxPkts = Range("D5").Value
    Dim srcIP As String: srcIP = Range("D6").Value
    Dim dstIP As String: dstIP = Range("D7").Value
    Dim proto As String: proto = Range("D8").Value
    Dim ipFlag As String: ipFlag = Range("D9").Value
    Dim startT As String: startT = Range("D10").Value
    Dim endT As String: endT = Range("D11").Value
    Dim debugFlag As String: debugFlag = Range("D12").Value
    

    cmd = """" & exePath & """"
    If filePath <> "" Then cmd = cmd & " -file """ & filePath & """"
    If outPath <> "" Then cmd = cmd & " -out """ & outPath & """"
    If maxPkts <> "" Then cmd = cmd & " -max " & maxPkts
    If srcIP <> "" Then cmd = cmd & " -source """ & srcIP & """"
    If dstIP <> "" Then cmd = cmd & " -destination """ & dstIP & """"
    If proto <> "" Then cmd = cmd & " -protocol """ & proto & """"
    If ipFlag <> "" Then cmd = cmd & " -IP """ & ipFlag & """"
    If startT <> "" Then cmd = cmd & " -startTime """ & startT & """"
    If endT <> "" Then cmd = cmd & " -endTime """ & endT & """"
    If debugFlag = "1" Then cmd = cmd & " -debug"
    
    Dim wsh As Object
    Set wsh = CreateObject("WScript.Shell")
    wsh.Run "cmd /k """ & cmd & """", 1, False
End Sub

Public Sub AddRunButton()
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Sheets("情報設定") ' 情報設定シートを指定
    
    Dim btn As Shape
    Set btn = ws.Shapes.AddFormControl(xlButtonControl, Range("D1").Left, Range("D1").Top, Range("D1").Width, Range("D1").Height) ' 位置調整
    btn.Name = "btnRun"
    btn.TextFrame.Characters.Text = "実行"
    btn.OnAction = "GenerateMermaidSequence"
End Sub

Sub CreateInfoSheet()
    Dim ws As Worksheet
    ' シートが存在するか確認
    On Error Resume Next
    Set ws = ThisWorkbook.Sheets("情報設定")
    On Error GoTo 0
    
    ' シートが存在しない場合のみ作成
    If ws Is Nothing Then
        Set ws = ThisWorkbook.Sheets.Add
        ws.Name = "情報設定"
    End If
    
    ' 既存のボタンを削除
    Dim btn As Shape ' btn変数を宣言
    For Each btn In ws.Shapes
        If btn.Name = "btnFileRef" Or btn.Name = "btnRun" Or btn.Name = "btnFolderRef" Then
            btn.Delete
        End If
    Next btn
    
    ' ラベルの追加
    ws.Range("C3").Value = "対象ファイル名:"
    ws.Range("C4").Value = "出力ファイル名:"
    ws.Range("C5").Value = "最大パケット数(全数なら何も入力しないでください):"
    ws.Range("C6").Value = "送信元IP(aaa.bbb.ccc.ddd):"
    ws.Range("C7").Value = "送信先IP(aaa.bbb.ccc.ddd):"
    ws.Range("C8").Value = "フィルタプロトコル名(例.DNS):"
    ws.Range("C9").Value = "送信元または送信先IP(aaa.bbb.ccc.ddd):"
    ws.Range("C10").Value = "フィルタ適用開始時刻(yyyy-mm-dd hh:mm:ss):"
    ws.Range("C11").Value = "フィルタ適用終了時刻(yyyy-mm-dd hh:mm:ss):"
    ws.Range("C12").Value = "デバッグモード(1:有効 0:無効):"

    ' 注意事項の追加
    ws.Range("C15").Value = "注意事項:"
    ws.Range("C16").Value = "このシートを使用する前に、必要な設定を確認してください。"
    ws.Range("C17").Value = "1. 対象ファイル名には、パケットキャプチャファイルのパスを指定してください。"
    ws.Range("C18").Value = "2. 出力ファイル名には、出力先のパスを指定してください。"
    ws.Range("C19").Value = "3. 最大パケット数は、必要に応じて指定してください。"
    ws.Range("C20").Value = "4. 送信元IPと送信先IPは、フィルタリングに使用されます。"
    ws.Range("C21").Value = "5. フィルタ適用時刻は、yyyy-mm-dd hh:mm:ss形式で指定してください。"
    ws.Range("C22").Value = "6. フィルタ適用終了時刻は、yyyy-mm-dd hh:mm:ss形式で指定してください。"
    ws.Range("C23").Value = "7. デバッグモードを有効にすると、詳細なログが生成されます。"
    ws.Range("C24").Value = "8. 設定が完了したら、実行ボタンを押してください。"
    ws.Range("C25").Value = "9. 出力ファイルは、指定したフォルダに保存されます。"
    ws.Range("C26").Value = "10. フィルタリング結果は、Mermaid形式で出力されます。"
    ws.Range("C27").Value = "11. マーメイド変換ボタンを押すことで、新しいシートが作成されシーケンスが作成されます。"

    ws.Range("C29").Value = "ファイル構成"
    ws.Range("C30").Value = "│ packet_sequence.exe"
    ws.Range("C31").Value = "├─execel_output"
    ws.Range("C32").Value = "│  ├─exec_command.bas"
    ws.Range("C33").Value = "│  ├─convert_mermaid.bas"
    ws.Range("C34").Value = "│  └─mermaid_output.xlsm(本ファイル)"
    ws.Range("C35").Value = "└config"
    ws.Range("C36").Value = "   └─config.pkseq"

    ws.Range("C38").Value = "基本的に大量にあるパケットファイルのシーケンスを作ることは試していないため、"
    ws.Range("C39").Value = "各フィルタ設定を行うことを推奨します。"
    ws.Range("C40").Value = "詳細な情報はwiresharkなどを使用して確認してください。"

    ws.Range("C42").Value = "各ページが消えた場合は、マクロからAuto_Openを実行してください。"
    
    ws.Range("C3:D12").Interior.Color = RGB(255, 255, 204) ' ラベルの背景色を変更
    ws.Range("C15:D42").Interior.Color = RGB(204, 255, 204) ' 注意事項の背景色を変更 - 正しいrange指定に修正

    ' ファイル参照ボタンの追加
    Dim btnFileRef As Shape
    Set btnFileRef = ws.Shapes.AddFormControl(xlButtonControl, Range("E3").Left, Range("E3").Top, Range("E3").Width, Range("E3").Height) ' 位置調整
    btnFileRef.Name = "btnFileRef"
    btnFileRef.TextFrame.Characters.Text = "ファイル参照"
    btnFileRef.OnAction = "BrowseFile"
    
    ' フォルダ参照ボタンの追加
    Dim btnFolderRef As Shape
    Set btnFolderRef = ws.Shapes.AddFormControl(xlButtonControl, Range("E4").Left, Range("E4").Top, Range("E4").Width, Range("E4").Height) ' 位置調整
    btnFolderRef.Name = "btnFolderRef"
    btnFolderRef.TextFrame.Characters.Text = "フォルダ参照"
    btnFolderRef.OnAction = "BrowseFolder"

    ' セル D12 のデータ検証設定：0か1のみ選択可能
    With ws.Range("D12").Validation
        .Delete
        .Add Type:=xlValidateList, AlertStyle:=xlValidAlertStop, Operator:=xlBetween, Formula1:="0,1"
        .IgnoreBlank = True
        .InCellDropdown = True
End With
    
    ' 実行ボタンの追加
    AddRunButton
End Sub

Sub BrowseFile()
    Dim fd As fileDialog
    Set fd = Application.fileDialog(msoFileDialogFilePicker)
    
    With fd
        .Title = "パケットキャプチャファイルを選択"
        .AllowMultiSelect = False
        .Filters.Clear
        .Filters.Add "PCAP Files", "*.pcap; *.pcapng; *.cap"
        If .Show = -1 Then
            ThisWorkbook.Sheets("情報設定").Range("D3").Value = .SelectedItems(1)
        Else
            msgBox "ファイルが選択されませんでした", vbExclamation
        End If
    End With
    
    Set fd = Nothing
End Sub

Sub BrowseFolder()
    Dim fd As FileDialog
    Set fd = Application.FileDialog(msoFileDialogFolderPicker)
    
    With fd
        .Title = "出力フォルダを選択"
        .AllowMultiSelect = False
        If .Show = -1 Then
            ThisWorkbook.Sheets("情報設定").Range("D4").Value = .SelectedItems(1) & "\output.md"
        Else
            MsgBox "フォルダが選択されませんでした", vbExclamation
        End If
    End With
    
    Set fd = Nothing
End Sub

Sub Auto_Open()
    CreateInfoSheet
    AddMermaidConversionButton ' マーメイド変換ボタンを追加
End Sub

