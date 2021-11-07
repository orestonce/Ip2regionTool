#ifndef TXTTODBFORM_H
#define TXTTODBFORM_H

#include <QWidget>

namespace Ui {
class TxtToDbForm;
}

class TxtToDbForm : public QWidget
{
    Q_OBJECT

public:
    explicit TxtToDbForm(QWidget *parent = 0);
    ~TxtToDbForm();
private slots:
    void on_pushButton_selectDb_clicked();

    void on_pushButton_selectTxt_clicked();

    void on_pushButton_startConvert_clicked();
private:
    void refresh_startConvert_Enable();
private:
    Ui::TxtToDbForm *ui;
};

#endif // TXTTODBFORM_H
