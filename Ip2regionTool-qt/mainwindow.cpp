#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QFileDialog>
#include <QMessageBox>
#include <QDebug>
#include "ip2region.h"

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    this->setConv_IsRuning(false);
    for(auto one : GetDbTypeList())
    {
        QString desc = QString::fromStdString(one.Desc);
        QVariant extName = QString::fromStdString(one.ExtName);

        ui->comboBox_fromType->addItem(desc, extName);
        if(one.SupportWrite)
        {
            ui->comboBox_toType->addItem(desc, extName);
        }
    }
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::setConv_IsRuning(bool runing)
{
    ui->lineEdit_fromName->setEnabled(runing == false);
    ui->lineEdit_toName->setEnabled(runing == false);
    ui->comboBox_fromType->setEnabled(runing == false);
    ui->comboBox_toType->setEnabled(runing == false);
    ui->pushButton_fromDir->setEnabled(runing == false);
    ui->pushButton_toDir->setEnabled(runing == false);
    ui->checkBox_VerifyFullUint32->setEnabled(runing == false);
    ui->checkBox_FillFullUint32->setEnabled(runing == false);
    ui->checkBox_MergeIpRange->setEnabled(runing == false);
    if (runing) {
        ui->pushButton_conv->setText("正在转换...");
    } else {
        ui->pushButton_conv->setText("开始转换");
    }
}

void MainWindow::on_pushButton_fromDir_clicked()
{
    QString extName = ui->comboBox_fromType->currentData().toString();
    QString input = QFileDialog::getOpenFileName(this, "",  "", "("+ extName + ")");
    if (!input.isEmpty()) {
        ui->lineEdit_fromName->setText(input);
    }
}

void MainWindow::on_pushButton_toDir_clicked()
{
    QString input = QFileDialog::getSaveFileName(this, "",  "", "");
    if (!input.isEmpty()) {
        ui->lineEdit_toName->setText(input);
    }
}

void MainWindow::on_pushButton_conv_clicked()
{
    ConvertDbReq req;
    req.FromName = ui->lineEdit_fromName->text().toStdString();
    req.ToName = ui->lineEdit_toName->text().toStdString();
    req.MergeIpRange = ui->checkBox_MergeIpRange->isChecked();
    req.FillFullUint32 = ui->checkBox_FillFullUint32->isChecked();
    req.VerifyFullUint32 = ui->checkBox_VerifyFullUint32->isChecked();
    if (req.FromName.empty() || req.ToName.empty()) {
        return;
    }
    this->setConv_IsRuning(true);

    req.FromType = ui->comboBox_fromType->currentText().toStdString();
    req.ToType = ui->comboBox_toType->currentText().toStdString();

    m_syncUi.AddRunFnOn_OtherThread([=](){
        std::string errMsg = ConvertDb(req);
        m_syncUi.AddRunFnOn_UiThread([=](){
            this->setConv_IsRuning(false);
            if (errMsg.empty()) {
                Toast::Instance()->SetSuccess("转换成功");
                return;
            }
            Toast::Instance()->SetError("失败 " + QString::fromStdString(errMsg));
        });
    });
}
