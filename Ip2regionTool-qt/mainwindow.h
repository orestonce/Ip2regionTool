#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();
private slots:
    void on_pushButton_input_db_clicked();

    void on_pushButton_output_txt_clicked();

    void on_pushButton_input_txt_clicked();

    void on_pushButton_output_db_clicked();

    void on_pushButton_DbToTxt_clicked();

    void on_pushButton_TxtToDb_clicked();

    void on_tabWidget_currentChanged(int index);

private:
    Ui::MainWindow *ui;
};

#endif // MAINWINDOW_H
