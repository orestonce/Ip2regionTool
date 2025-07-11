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
    void setConv_IsRuning(bool runing);
private slots:
    void on_pushButton_fromDir_clicked();

    void on_pushButton_toDir_clicked();

    void on_pushButton_conv_clicked();

private:
    Ui::MainWindow *ui;
    RunOnUiThread m_syncUi;
};

#endif // MAINWINDOW_H
