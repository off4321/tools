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
    ws.Range("C5").Value = "最大パケット数:"
    ws.Range("C6").Value = "送信元IP(aaa.bbb.ccc.ddd):"
    ws.Range("C7").Value = "送信先IP(aaa.bbb.ccc.ddd):"
    ws.Range("C8").Value = "フィルタプロトコル名(例.DNS):"
    ws.Range("C9").Value = "送信元または送信先IP(aaa.bbb.ccc.ddd):"
    ws.Range("C10").Value = "フィルタ適用開始時刻(yyyy-mm-dd hh:mm:ss):"
    ws.Range("C11").Value = "フィルタ適用終了時刻(yyyy-mm-dd hh:mm:ss):"
    
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
End Sub

