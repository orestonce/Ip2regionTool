#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include "ip2region.h"

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();
    void setTxtToXdb_IsRuning(bool runing);
private slots:
    void on_pushButton_input_db_clicked();

    void on_pushButton_output_txt_clicked();

    void on_pushButton_input_txt_clicked();

    void on_pushButton_output_db_clicked();

    void on_pushButton_DbToTxt_clicked();

    void on_pushButton_TxtToDb_clicked();

    void on_tabWidget_currentChanged(int index);

    void on_pushButton_input_regin_csv_clicked();

    void on_pushButton_xdb_srcFile_clicked();

    void on_pushButton_xdb_dstFile_clicked();

    void on_pushButton_xdb_start_clicked();

private:
    Ui::MainWindow *ui;
    RunOnUiThread m_syncUi;
};

#endif // MAINWINDOW_H
