#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include "dbtotxtform.h"
#include "txttodbform.h"

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
    void on_actionDbToTxt_triggered();

    void on_actionTxtToDb_triggered();
private:
    void refresh_status(int idx);
private:
    Ui::MainWindow *ui;
    DbToTxtForm* m_form1;
    TxtToDbForm* m_form2;
};

#endif // MAINWINDOW_H
